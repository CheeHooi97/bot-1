package main

import (
	"bot-1/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/gorilla/websocket"
)

// Candle represents a Binance kline/candle message
type Candle struct {
	OpenTime  int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	CloseTime int64
	// other fields ignored
}

// Strategy parameters
const (
	rsiLength      = 14
	volumeLookback = 10
	bbLength       = 20
	bbMult         = 2.0
)

// State to track positions: 0=flat, 1=long, -1=short
var state = 0

// Store closes and volumes for indicators
var closes []float64
var volumes []float64

// Binance client
var client *binance.Client

func main() {
	// load config
	config.LoadConfig()

	log.Println("Key:", config.BinanceApiKey)
	log.Println("Secret:", config.BinanceApiSecret)

	client = binance.NewClient(config.BinanceApiKey, config.BinanceApiSecret)

	symbol := "btcusdt"
	interval := "1m"

	// Connect to Binance WebSocket for kline data
	wsURL := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@kline_%s", symbol, interval)
	log.Printf("Connecting to %s", wsURL)

	c, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{})
	if err != nil {
		log.Fatal("WebSocket dial error:", err)
	}
	defer c.Close()

	// Handle graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Println("test")
	log.Println("Bot started. Waiting for live candle data...")

	for {
		select {
		case <-interrupt:
			log.Println("Received interrupt, shutting down...")
			return
		default:
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				time.Sleep(time.Second * 3)
				continue
			}

			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("JSON unmarshal error:", err)
				continue
			}

			kline, ok := msg["k"].(map[string]interface{})
			if !ok {
				continue
			}

			isFinal := kline["x"].(bool)
			if !isFinal {
				// Only process closed candles
				continue
			}

			candle, err := parseCandle(kline)
			if err != nil {
				log.Println("Error parsing candle:", err)
				continue
			}

			log.Println("Open")
			log.Println(candle.Open)

			log.Println("Close")
			log.Println(candle.Close)

			log.Println("High")
			log.Println(candle.High)

			log.Println("Low")
			log.Println(candle.Low)

			// processCandle(candle, symbol)
		}
	}
}

// parseCandle extracts candle info from Binance kline JSON
func parseCandle(k map[string]interface{}) (Candle, error) {
	var c Candle
	var err error

	c.OpenTime = int64(k["t"].(float64))
	c.Open, err = strconv.ParseFloat(k["o"].(string), 64)
	if err != nil {
		return c, err
	}
	c.High, err = strconv.ParseFloat(k["h"].(string), 64)
	if err != nil {
		return c, err
	}
	c.Low, err = strconv.ParseFloat(k["l"].(string), 64)
	if err != nil {
		return c, err
	}
	c.Close, err = strconv.ParseFloat(k["c"].(string), 64)
	if err != nil {
		return c, err
	}
	c.Volume, err = strconv.ParseFloat(k["v"].(string), 64)
	if err != nil {
		return c, err
	}
	c.CloseTime = int64(k["T"].(float64))

	return c, nil
}

// processCandle runs the strategy logic on each closed candle
func processCandle(c Candle, symbol string) {
	// Append closes and volumes
	closes = append(closes, c.Close)
	volumes = append(volumes, c.Volume)

	// Maintain max length for indicators
	if len(closes) > 200 {
		closes = closes[1:]
		volumes = volumes[1:]
	}

	if len(closes) < rsiLength || len(volumes) < volumeLookback || len(closes) < bbLength {
		// Not enough data yet
		return
	}

	rsiVal := calcRSI(closes, rsiLength)
	avgVolume := sma(volumes, volumeLookback)
	highVolume := c.Volume > avgVolume*1.5
	extremeHighVolume := c.Volume > avgVolume*2

	// Candle color
	greenCandle := c.Close > c.Open
	redCandle := c.Close < c.Open

	// Wick calculations
	highLowDiff := c.High - c.Low
	if highLowDiff == 0 {
		return // avoid division by zero
	}
	topWick := c.High - math.Max(c.Open, c.Close)
	bottomWick := math.Min(c.Open, c.Close) - c.Low
	topWickPerc := (topWick / highLowDiff) * 100
	bottomWickPerc := (bottomWick / highLowDiff) * 100

	// Bollinger Bands
	basis := sma(closes[len(closes)-bbLength:], bbLength)
	stdDev := stddev(closes[len(closes)-bbLength:], basis)
	upper := basis + bbMult*stdDev
	lower := basis - bbMult*stdDev

	rawBuy := (rsiVal < 35 && highVolume && (greenCandle || (redCandle && bottomWickPerc > 60))) || extremeHighVolume && c.Close <= lower
	rawSell := (rsiVal > 65 && highVolume && (redCandle || (greenCandle && topWickPerc > 60))) || extremeHighVolume && c.Close >= upper

	canLong := true  // Assuming both long and short allowed
	canShort := true // Adjust as needed

	buySignal := rawBuy && canLong
	sellSignal := rawSell && canShort

	log.Printf("Candle Close: %.2f RSI: %.2f Volume: %.2f AvgVol: %.2f Buy: %v Sell: %v", c.Close, rsiVal, c.Volume, avgVolume, buySignal, sellSignal)

	if buySignal && (state == 0 || state == -1) {
		log.Println("BUY signal triggered")
		// err := placeOrder(symbol, "BUY")
		// if err == nil {
		// 	state = 1
		// }
	}

	if sellSignal && (state == 0 || state == 1) {
		log.Println("SELL signal triggered")
		// err := placeOrder(symbol, "SELL")
		// if err == nil {
		// 	state = -1
		// }
	}
}

// placeOrder submits a market order to Binance
func placeOrder(symbol, side string) error {
	sideType := binance.SideTypeBuy
	if side == "SELL" {
		sideType = binance.SideTypeSell
	}

	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(sideType).
		Type(binance.OrderTypeMarket).
		Quantity("0.001"). // example qty, adjust as needed
		Do(context.Background())
	if err != nil {
		log.Println("Order failed:", err)
		return err
	}

	log.Printf("%s order placed: %+v\n", side, order)
	return nil
}

// RSI calculation
func calcRSI(data []float64, length int) float64 {
	if len(data) < length+1 {
		return 50.0 // neutral
	}
	var gainSum, lossSum float64

	for i := len(data) - length; i < len(data); i++ {
		diff := data[i] - data[i-1]
		if diff > 0 {
			gainSum += diff
		} else {
			lossSum -= diff
		}
	}
	if lossSum == 0 {
		return 100
	}
	rs := gainSum / lossSum
	rsi := 100 - (100 / (1 + rs))
	return rsi
}

// Simple moving average
func sma(data []float64, length int) float64 {
	if len(data) < length {
		return 0
	}
	sum := 0.0
	for i := len(data) - length; i < len(data); i++ {
		sum += data[i]
	}
	return sum / float64(length)
}

// Standard deviation
func stddev(data []float64, mean float64) float64 {
	variance := 0.0
	for _, v := range data {
		diff := v - mean
		variance += diff * diff
	}
	return math.Sqrt(variance / float64(len(data)))
}
