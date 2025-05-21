package main

import (
	"bot-1/config"
	"bot-1/constant"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/gorilla/websocket"
)

type Candle struct {
	OpenTime  int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	CloseTime int64
}

const (
	rsiLength      = 14
	volumeLookback = 10
	bbLength       = 20
	bbMult         = 2.0
)

type SymbolData struct {
	mu      sync.Mutex
	Closes  []float64
	Volumes []float64
	State   int
}

var client *binance.Client

func main() {
	apiKey := config.BinanceApiKey
	apiSecret := config.BinanceApiSecret
	if apiKey == "" || apiSecret == "" {
		log.Fatal("Set BINANCE_API_KEY and BINANCE_API_SECRET environment variables")
	}
	client = binance.NewClient(apiKey, apiSecret)

	symbols := []string{"btcusdt", "ethusdt", "solusdt"}
	interval := "15m"

	var wg sync.WaitGroup
	for _, symbol := range symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			runBot(sym, interval)
		}(symbol)
	}

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutdown signal received. Exiting...")
}

func runBot(symbol, interval string) {
	data := &SymbolData{}
	wsURL := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@kline_%s", symbol, interval)
	log.Printf("[%s] Connecting to %s", symbol, wsURL)

	c, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{})
	if err != nil {
		log.Fatalf("[%s] WebSocket dial error: %v", symbol, err)
		return
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("[%s] Read error: %v", symbol, err)
			time.Sleep(3 * time.Second)
			continue
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("[%s] JSON unmarshal error: %v", symbol, err)
			continue
		}

		kline, ok := msg["k"].(map[string]interface{})
		if !ok || !kline["x"].(bool) {
			continue
		}

		candle, err := parseCandle(kline)
		if err != nil {
			log.Printf("[%s] Error parsing candle: %v", symbol, err)
			continue
		}

		processCandle(candle, symbol, data)
	}
}

func parseCandle(k map[string]interface{}) (Candle, error) {
	var c Candle
	var err error

	c.OpenTime = int64(k["t"].(float64))
	c.Open, _ = strconv.ParseFloat(k["o"].(string), 64)
	c.High, _ = strconv.ParseFloat(k["h"].(string), 64)
	c.Low, _ = strconv.ParseFloat(k["l"].(string), 64)
	c.Close, _ = strconv.ParseFloat(k["c"].(string), 64)
	c.Volume, err = strconv.ParseFloat(k["v"].(string), 64)
	c.CloseTime = int64(k["T"].(float64))
	return c, err
}

func processCandle(c Candle, symbol string, data *SymbolData) {
	data.mu.Lock()
	defer data.mu.Unlock()

	data.Closes = append(data.Closes, c.Close)
	data.Volumes = append(data.Volumes, c.Volume)

	if len(data.Closes) > 200 {
		data.Closes = data.Closes[1:]
		data.Volumes = data.Volumes[1:]
	}

	if len(data.Closes) < rsiLength || len(data.Volumes) < volumeLookback || len(data.Closes) < bbLength {
		return
	}

	rsiVal := calcRSI(data.Closes, rsiLength)
	avgVolume := sma(data.Volumes, volumeLookback)
	highVolume := c.Volume > avgVolume*1.5
	extremeHighVolume := c.Volume > avgVolume*2
	greenCandle := c.Close > c.Open
	redCandle := c.Close < c.Open

	highLowDiff := c.High - c.Low
	if highLowDiff == 0 {
		return
	}
	topWick := c.High - math.Max(c.Open, c.Close)
	bottomWick := math.Min(c.Open, c.Close) - c.Low
	topWickPerc := (topWick / highLowDiff) * 100
	bottomWickPerc := (bottomWick / highLowDiff) * 100

	basis := sma(data.Closes[len(data.Closes)-bbLength:], bbLength)
	stdDev := stddev(data.Closes[len(data.Closes)-bbLength:], basis)
	upper := basis + bbMult*stdDev
	lower := basis - bbMult*stdDev

	rawBuy := (rsiVal < 35 && highVolume && (greenCandle || (redCandle && bottomWickPerc > 60))) || extremeHighVolume && c.Close <= lower
	rawSell := (rsiVal > 65 && highVolume && (redCandle || (greenCandle && topWickPerc > 60))) || extremeHighVolume && c.Close >= upper

	buySignal := rawBuy && (data.State == 0 || data.State == -1)
	sellSignal := rawSell && (data.State == 0 || data.State == 1)

	log.Printf("[%s] Close: %.2f RSI: %.2f Vol: %.2f AvgVol: %.2f Buy: %v Sell: %v", symbol, c.Close, rsiVal, c.Volume, avgVolume, buySignal, sellSignal)

	if buySignal {
		log.Printf("[%s] BUY signal", symbol)
		if err := placeOrder(symbol, "BUY"); err == nil {
			data.State = 1
		}
	}

	if sellSignal {
		log.Printf("[%s] SELL signal", symbol)
		if err := placeOrder(symbol, "SELL"); err == nil {
			data.State = -1
		}
	}
}

func placeOrder(symbol, side string) error {
	sideType := binance.SideTypeBuy
	if side == "SELL" {
		sideType = binance.SideTypeSell
	}

	quantity, ok := constant.QuantityMap[symbol]
	if !ok {
		log.Printf("[%s] No quantity defined, skipping order", symbol)
		return fmt.Errorf("no quantity set for %s", symbol)
	}

	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(sideType).
		Type(binance.OrderTypeMarket).
		Quantity(quantity).
		Do(context.Background())
	if err != nil {
		log.Printf("[%s] Order failed: %v", symbol, err)
		return err
	}

	log.Printf("[%s] %s order placed: %+v", symbol, side, order)
	return nil
}

func calcRSI(data []float64, length int) float64 {
	if len(data) < length+1 {
		return 50.0
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
	return 100 - (100 / (1 + rs))
}

func sma(data []float64, length int) float64 {
	sum := 0.0
	for i := len(data) - length; i < len(data); i++ {
		sum += data[i]
	}
	return sum / float64(length)
}

func stddev(data []float64, mean float64) float64 {
	variance := 0.0
	for _, v := range data {
		diff := v - mean
		variance += diff * diff
	}
	return math.Sqrt(variance / float64(len(data)))
}
