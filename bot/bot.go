package bot

import (
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
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	CloseTime time.Time
	IsFinal   bool
}

var (
	closes         []float64
	volumes        []float64
	rsiLength      = 14
	volumeLookback = 20
	bbLength       = 20
	bbMult         = 2.0
	tradeUSDT      = 100.0
	balance        = 1000.0
	positionSize   float64
	entryPrice     float64
	state          = 0 // 0 = neutral, 1 = long, -1 = short
	client         *binance.Client
)

func Bot(symbol, interval string) {
	// Step 1: Fetch historical candles
	history, err := fetchHistoricalCandles(symbol, interval, 200)
	if err != nil {
		log.Fatal("Error fetching historical candles:", err)
	}
	for _, c := range history {
		processCandle(c, symbol)
	}

	// Step 2: Start WebSocket
	go startWebSocket(symbol, interval)
	waitForShutdown()
}

func fetchHistoricalCandles(symbol, interval string, limit int) ([]Candle, error) {
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/klines?symbol=%s&interval=%s&limit=%d", symbol, interval, limit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var candles []Candle
	for _, item := range data {
		open := parseStringToFloat(item[1])
		high := parseStringToFloat(item[2])
		low := parseStringToFloat(item[3])
		close := parseStringToFloat(item[4])
		volume := parseStringToFloat(item[5])
		closeTime := time.UnixMilli(int64(item[6].(float64)))

		candles = append(candles, Candle{
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
			CloseTime: closeTime,
			IsFinal:   true,
		})
	}
	return candles, nil
}

func parseStringToFloat(s interface{}) float64 {
	val, _ := strconv.ParseFloat(s.(string), 64)
	return val
}

func startWebSocket(symbol, interval string) {
	url := fmt.Sprintf("wss://fstream.binance.com/ws/%s@kline_%s", symbol, interval)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("WebSocket dial error:", err)
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal(message, &raw); err != nil {
			continue
		}

		kline, ok := raw["k"].(map[string]interface{})
		if !ok {
			continue
		}

		if !kline["x"].(bool) {
			continue
		}

		candle := Candle{
			Open:      parseStringToFloat(kline["o"]),
			High:      parseStringToFloat(kline["h"]),
			Low:       parseStringToFloat(kline["l"]),
			Close:     parseStringToFloat(kline["c"]),
			Volume:    parseStringToFloat(kline["v"]),
			CloseTime: time.UnixMilli(int64(kline["T"].(float64))),
			IsFinal:   true,
		}

		processCandle(candle, symbol)
	}
}

func waitForShutdown() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	log.Println("Shutting down.")
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

	rawBuy := (rsiVal < 35 && highVolume && (greenCandle || (redCandle && bottomWickPerc > 60))) || (extremeHighVolume && c.Close <= lower)
	rawSell := (rsiVal > 65 && highVolume && (redCandle || (greenCandle && topWickPerc > 60))) || (extremeHighVolume && c.Close >= upper)

	canLong := true  // Assuming both long and short allowed
	canShort := true // Adjust as needed

	buySignal := rawBuy && canLong
	sellSignal := rawSell && canShort

	// === ONLY LOG ON POSITION CHANGE, NOT EVERY CANDLE ===

	// Closing SHORT position and opening LONG
	if buySignal && (state == 0 || state == -1) {
		if state == -1 {
			// Closing short position first: buy BTC to cover short
			closeAmount := math.Abs(positionSize)
			profit := (entryPrice - c.Close) * closeAmount
			balance += tradeUSDT + profit
			log.Printf("Closing SHORT position: bought %.4f BTC at %.2f, profit: %.2f USDT, balance: %.2f USDT", closeAmount, c.Close, profit, balance)
			positionSize = 0
			entryPrice = 0
			state = 0
		}
		if balance >= tradeUSDT {
			// Open long position
			size := tradeUSDT / c.Close
			positionSize = size
			entryPrice = c.Close
			balance -= tradeUSDT
			state = 1
			log.Printf("Opened LONG position: bought %.4f BTC at %.2f, spent %.2f USDT, remaining balance %.2f USDT", size, c.Close, tradeUSDT, balance)
		} else {
			log.Println("Insufficient balance to open LONG position")
		}
		return
	}

	// Closing LONG position and opening SHORT
	if sellSignal && (state == 0 || state == 1) {
		if state == 1 {
			// Closing long position first: sell BTC
			profit := (c.Close - entryPrice) * positionSize
			balance += tradeUSDT + profit // Return initial tradeUSDT + profit
			log.Printf("Closing LONG position: sold %.4f BTC at %.2f, profit: %.2f USDT, balance: %.2f USDT", positionSize, c.Close, profit, balance)
			positionSize = 0
			entryPrice = 0
			state = 0
		}
		if balance >= tradeUSDT {
			// Open short position
			size := tradeUSDT / c.Close
			positionSize = -size
			entryPrice = c.Close
			balance -= tradeUSDT
			state = -1
			log.Printf("Opened SHORT position: sold short %.4f BTC at %.2f, used %.2f USDT margin, remaining balance %.2f USDT", size, c.Close, tradeUSDT, balance)
		} else {
			log.Println("Insufficient balance to open SHORT position")
		}
		return
	}

	// If holding a position (long or short), no logs to avoid flooding logs every candle
}

// placeOrder submits a market order to Binance (currently unused)
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
