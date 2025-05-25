package bot

import (
	"bot-1/config"
	"bot-1/constant"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	closes          []float64
	volumes         []float64
	rsiLength       = 14
	volumeLookback  = 20
	bbLength        = 20
	bbMult          = 2.0
	tradeUSDT       = 500.0
	balance         = 10000.0
	totalProfitLoss = 0.0
	positionSize    float64
	entryPrice      float64
	state           = 0 // 0 = neutral, 1 = long, -1 = short
	client          *binance.Client
	stopLossPercent float64
)

// Bot runs the trading bot on given symbol, interval and stop loss percent
func Bot(symbol, interval, token string, slPercent float64) {
	stopLossPercent = slPercent
	client = binance.NewClient("", "")

	// Fetch historical candles
	history, err := fetchHistoricalCandles(strings.ToUpper(symbol), interval)
	if err != nil {
		log.Fatal("Error fetching historical candles:", err)
	}
	for _, c := range history {
		closes = append(closes, c.Close)
		volumes = append(volumes, c.Volume)
	}

	// Keep buffer size trimmed
	if len(closes) > 500 {
		closes = closes[len(closes)-500:]
		volumes = volumes[len(volumes)-500:]
	}

	// Start WebSocket
	go startWebSocket(strings.ToLower(symbol), interval, token)
	waitForShutdown()
}

func fetchHistoricalCandles(symbol, interval string) ([]Candle, error) {
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/klines?symbol=%s&interval=%s", symbol, interval)
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

func startWebSocket(symbol, interval, token string) {
	url := fmt.Sprintf("wss://fstream.binance.com/ws/%s@kline_%s", symbol, interval)
	log.Println("Connecting to ", url)
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

		if !kline["x"].(bool) { // only closed candles
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

		processCandle(candle, symbol, token)
	}
}

func waitForShutdown() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	log.Println("Shutting down.")
}

// processCandle runs the strategy logic on each closed candle
func processCandle(c Candle, symbol, token string) {
	s := strings.ToUpper(symbol[:len(symbol)-4])

	closes = append(closes, c.Close)
	volumes = append(volumes, c.Volume)

	if len(closes) > 500 {
		closes = closes[1:]
		volumes = volumes[1:]
	}

	if len(closes) < rsiLength || len(volumes) < volumeLookback || len(closes) < bbLength {
		return
	}

	step := constant.StepMap[symbol]

	pricePrecision := constant.SymbolPrecisionMap[symbol][0]
	amountPrecision := constant.SymbolPrecisionMap[symbol][1]

	price := fmt.Sprintf("%%.%df", pricePrecision)
	amount := fmt.Sprintf("%%.%df", amountPrecision)

	rsiVal := calcRSI(closes, rsiLength)
	avgVolume := sma(volumes, volumeLookback)
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

	basis := sma(closes[len(closes)-bbLength:], bbLength)
	stdDev := stddev(closes[len(closes)-bbLength:], basis)
	upper := basis + bbMult*stdDev
	lower := basis - bbMult*stdDev

	rawBuy := (rsiVal < 35 && highVolume && (greenCandle || (redCandle && bottomWickPerc > 60))) || (extremeHighVolume && c.Close <= lower)
	rawSell := (rsiVal > 65 && highVolume && (redCandle || (greenCandle && topWickPerc > 60))) || (extremeHighVolume && c.Close >= upper)

	buySignal := rawBuy
	sellSignal := rawSell

	// === STOP LOSS CHECK ===
	if state == 1 && c.Close <= entryPrice*(1-stopLossPercent/100) {
		profit := (c.Close - entryPrice) * positionSize
		percentChange := ((c.Close - entryPrice) / entryPrice) * 100
		balance += tradeUSDT + profit
		// placeOrder(symbol, "SELL", positionSize)
		a := fmt.Sprintf("STOP LOSS [LONG]\nAmount: "+amount+"\nPrice: "+price+"\nPercent changed: %.2f\nLoss: %.2f USDT\nBalance: %.2f USDT", positionSize, s, c.Close, percentChange, profit, balance)
		log.Println(a)
		sendTelegramMessage(token, a)
		state = 0
		positionSize = 0
		entryPrice = 0
		totalProfitLoss += profit
		b := fmt.Sprintf("Total profit/loss : %.2f", totalProfitLoss)
		log.Println(b)
		sendTelegramMessage(token, b)
		return
	}
	if state == -1 && c.Close >= entryPrice*(1+stopLossPercent/100) {
		closeAmount := math.Abs(positionSize)
		profit := (entryPrice - c.Close) * closeAmount
		percentChange := ((c.Close - entryPrice) / entryPrice) * 100
		balance += tradeUSDT + profit
		// placeOrder(symbol, "BUY", positionSize)
		a := fmt.Sprintf("STOP LOSS [SHORT]\nAmount: "+amount+"\nPrice: "+price+"\nPercent changed: %.2f\nLoss: %.2f USDT\nBalance: %.2f USDT", closeAmount, s, c.Close, percentChange, profit, balance)
		log.Println(a)
		sendTelegramMessage(token, a)
		state = 0
		positionSize = 0
		entryPrice = 0
		totalProfitLoss += profit
		b := fmt.Sprintf("Total profit/loss : %.2f", totalProfitLoss)
		log.Println(b)
		sendTelegramMessage(token, b)
		return
	}

	// === TRADING LOGIC ===
	if state == 0 {
		// Neutral: open position on any signal
		if buySignal {
			if balance >= tradeUSDT {
				size := roundUpToStep(tradeUSDT/c.Close, step)
				positionSize = size
				entryPrice = c.Close
				balance -= tradeUSDT
				state = 1
				// placeOrder(symbol, "BUY", positionSize)
				a := fmt.Sprintf("[LONG]\nAmount: "+amount+"\nPrice: "+price+"\nStop loss: "+price+"\nBalance: %.2f", size, s, c.Close, c.Close*(1-stopLossPercent/100), balance)
				log.Println(a)
				sendTelegramMessage(token, a)
			} else {
				a := "Insufficient balance to open LONG position"
				log.Println(a)
				sendTelegramMessage(token, a)
			}
			return
		}
		if sellSignal {
			if balance >= tradeUSDT {
				size := roundUpToStep(tradeUSDT/c.Close, step)
				positionSize = -size
				entryPrice = c.Close
				balance -= tradeUSDT
				state = -1
				// placeOrder(symbol, "SELL", positionSize)
				a := fmt.Sprintf("[SHORT]\nAmount: "+amount+"\nPrice: "+price+"\nStop loss: "+price+"\nBalance: %.2f", size, s, c.Close, c.Close*(1+stopLossPercent/100), balance)
				log.Println(a)
				sendTelegramMessage(token, a)
			} else {
				a := "Insufficient balance to open SHORT position"
				log.Println(a)
				sendTelegramMessage(token, a)
			}
			return
		}
	} else if state == 1 {
		// Long position: close only on sell signal
		if sellSignal {
			profit := (c.Close - entryPrice) * positionSize
			percentChange := ((c.Close - entryPrice) / entryPrice) * 100
			balance += tradeUSDT + profit
			// placeOrder(symbol, "SELL", positionSize)
			a := fmt.Sprintf("Closed [LONG]\nAmount: "+amount+"\nPrice: "+price+"\nPercent changed: %.2f\nProfit: %.2f USDT\nBalance: %.2f USDT", positionSize, s, c.Close, percentChange, profit, balance)
			log.Println(a)
			sendTelegramMessage(token, a)
			state = 0
			positionSize = 0
			entryPrice = 0
			totalProfitLoss += profit
			b := fmt.Sprintf("Total profit/loss : %.2f", totalProfitLoss)
			log.Println(b)
			sendTelegramMessage(token, b)
			return
		}
	} else if state == -1 {
		// Short position: close only on buy signal
		if buySignal {
			closeAmount := math.Abs(positionSize)
			profit := (entryPrice - c.Close) * closeAmount
			percentChange := ((c.Close - entryPrice) / entryPrice) * 100
			balance += tradeUSDT + profit
			// placeOrder(symbol, "BUY", positionSize)
			a := fmt.Sprintf("Closed [SHORT]\nAmount: "+amount+"\nPrice: "+price+"\nPercent changed: %.2f\nProfit: %.2f USDT\nBalance: %.2f USDT", closeAmount, s, c.Close, percentChange, profit, balance)
			log.Println(a)
			sendTelegramMessage(token, a)
			state = 0
			positionSize = 0
			entryPrice = 0
			totalProfitLoss += profit
			b := fmt.Sprintf("Total profit/loss : %.2f", totalProfitLoss)
			log.Println(b)
			sendTelegramMessage(token, b)
			return
		}
	}
}

func placeOrder(symbol string, side string, quantity float64) {
	endpoint := "https://fapi.binance.com/fapi/v1/order"
	timestamp := time.Now().UnixMilli()

	// Convert quantity to string with 4 decimals
	qtyStr := strconv.FormatFloat(quantity, 'f', 4, 64)

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side) // "BUY" or "SELL"
	params.Set("type", "MARKET")
	params.Set("quantity", qtyStr)
	params.Set("timestamp", strconv.FormatInt(timestamp, 10))

	// Sign the query
	queryString := params.Encode()
	signature := sign(queryString, config.BinanceApiSecret)
	params.Set("signature", signature)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(params.Encode()))
	if err != nil {
		log.Println("Error creating order request:", err)
		return
	}

	req.Header.Set("X-MBX-APIKEY", config.BinanceApiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error placing order:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("Order response:", string(body))
}

// sign generates HMAC-SHA256 signature
func sign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func calcRSI(closes []float64, length int) float64 {
	if len(closes) < length+1 {
		return 0
	}
	var gains, losses float64
	for i := len(closes) - length; i < len(closes); i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}
	if losses == 0 {
		return 100
	}
	rs := gains / losses
	rsi := 100 - (100 / (1 + rs))
	return rsi
}

func sma(data []float64, length int) float64 {
	if len(data) < length {
		return 0
	}
	sum := 0.0
	for _, v := range data[len(data)-length:] {
		sum += v
	}
	return sum / float64(length)
}

func stddev(data []float64, mean float64) float64 {
	var sum float64
	for _, v := range data {
		sum += (v - mean) * (v - mean)
	}
	return math.Sqrt(sum / float64(len(data)))
}

func roundUpToStep(value, step float64) float64 {
	return math.Ceil(value/step) * step
}
