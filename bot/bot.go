package bot

import (
	"archive/zip"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"dca-bot/config"
	"dca-bot/constant"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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
	balance         = 10000.0
	totalProfitLoss = 0.0
	entryPrice      float64
	state           = 0 // 0 = neutral, 1 = LONG, -1 = SHORT
	stopLossPercent float64
	numOfWin        = 0
	numOfLose       = 0
)

// Bot runs the trading bot on given symbol, interval and stop loss percent
func Bot(symbol, interval, token string, slPercent float64) {
	stopLossPercent = slPercent

	history, err := loadHistoricalFromVision(strings.ToUpper(symbol), interval, 500)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range history {
		closes = append(closes, c.Close)
		volumes = append(volumes, c.Volume)
	}

	sendTelegramMessage(token, fmt.Sprintf("%s %s started", symbol, interval))
	go startWebSocket(strings.ToLower(symbol), interval, token)
	waitForShutdown()
}

func loadHistoricalFromVision(symbol, interval string, limit int) ([]Candle, error) {
	var all []Candle

	now := time.Now()
	for i := 0; len(all) < limit && i < 6; i++ {
		t := now.AddDate(0, -i, 0)
		ym := t.Format("2006-01")

		candles, err := loadVisionMonth(symbol, interval, ym)
		if err == nil {
			all = append(candles, all...)
		}
	}

	if len(all) < limit {
		return nil, fmt.Errorf("not enough historical candles")
	}

	return all[len(all)-limit:], nil
}

func loadVisionMonth(symbol, interval, ym string) ([]Candle, error) {
	url := fmt.Sprintf(
		"https://data.binance.vision/data/futures/um/monthly/klines/%s/%s/%s-%s-%s.zip",
		symbol, interval, symbol, interval, ym,
	)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("vision download failed")
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	zr, _ := zip.NewReader(bytes.NewReader(data), int64(len(data)))

	var candles []Candle
	for _, f := range zr.File {
		rc, _ := f.Open()
		r := csv.NewReader(rc)

		for {
			row, err := r.Read()
			if err != nil {
				break
			}

			candles = append(candles, Candle{
				Open:      parseFloat(row[1]),
				High:      parseFloat(row[2]),
				Low:       parseFloat(row[3]),
				Close:     parseFloat(row[4]),
				Volume:    parseFloat(row[5]),
				CloseTime: time.UnixMilli(parseInt(row[6])),
				IsFinal:   true,
			})
		}
		rc.Close()
	}
	return candles, nil
}

func openOrderSide(state int) string {
	if state == 1 {
		return "BUY"
	}
	if state == -1 {
		return "SELL"
	}
	return ""
}

func closeOrderSide(state int) string {
	if state == 1 {
		return "SELL"
	}
	if state == -1 {
		return "BUY"
	}
	return ""
}

func positionLabel(state int) string {
	if state == 1 {
		return "LONG"
	}
	if state == -1 {
		return "SHORT"
	}
	return "NONE"
}

func parseFloat(v any) float64 {
	f, _ := strconv.ParseFloat(v.(string), 64)
	return f
}

func parseInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func parseStringToFloat(s any) float64 {
	val, _ := strconv.ParseFloat(s.(string), 64)
	return val
}

func startWebSocket(symbol, interval, token string) {
	wsURL := fmt.Sprintf("wss://fstream.binance.com/ws/%s@kline_%s", symbol, interval)

	for {
		c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				c.Close()
				break
			}

			var raw map[string]any
			json.Unmarshal(msg, &raw)

			k, ok := raw["k"].(map[string]any)
			if !ok || !k["x"].(bool) {
				continue
			}

			candle := Candle{
				Open:   parse(k["o"]),
				High:   parse(k["h"]),
				Low:    parse(k["l"]),
				Close:  parse(k["c"]),
				Volume: parse(k["v"]),
			}

			processCandle(candle, symbol, token)
		}
	}
}

func waitForShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
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

	if len(closes) < 20 {
		return
	}

	rsi := calcRSI(closes, rsiLength)
	avgVol := sma(volumes, volumeLookback)

	basis := sma(closes[len(closes)-bbLength:], bbLength)
	std := stddev(closes[len(closes)-bbLength:], basis)
	upper := basis + bbMult*std
	lower := basis - bbMult*std

	buySignal := rsi < 35 && c.Volume > avgVol*1.5 && c.Close <= lower
	sellSignal := rsi > 65 && c.Volume > avgVol*1.5 && c.Close >= upper

	size := constant.QuantityMap[symbol]
	positionSize := strconv.FormatFloat(size, 'f', constant.SymbolPrecisionMap[symbol][1], 64)
	price := strconv.FormatFloat(c.Close, 'f', constant.SymbolPrecisionMap[symbol][0], 64)

	/* ===== STOP LOSS ===== */

	if state != 0 {
		stop := entryPrice * (1 - stopLossPercent/100)
		if state == -1 {
			stop = entryPrice * (1 + stopLossPercent/100)
		}

		if (state == 1 && c.Close <= stop) || (state == -1 && c.Close >= stop) {
			side := closeOrderSide(state)
			placeOrder(symbol, side)

			profit := (c.Close - entryPrice) * size
			if state == -1 {
				profit = (entryPrice - c.Close) * size
			}

			balance += size*c.Close + profit
			totalProfitLoss += profit
			numOfLose++

			sendTelegramMessage(token,
				fmt.Sprintf("STOP LOSS [%s]\nPrice: %s\nPNL: %.2f\nBalance: %.2f",
					positionLabel(state), price, profit, balance))

			state, entryPrice = 0, 0
			return
		}
	}

	/* ===== OPEN ===== */

	if state == 0 {
		if buySignal {
			state = 1
		} else if sellSignal {
			state = -1
		} else {
			return
		}

		entryPrice = c.Close
		balance -= size * c.Close
		placeOrder(symbol, openOrderSide(state))

		sendTelegramMessage(token,
			fmt.Sprintf("OPEN [%s]\nAmount: %s %s\nPrice: %s\nBalance: %.2f",
				positionLabel(state), positionSize, s, price, balance))
		return
	}

	/* ===== CLOSE ===== */

	if (state == 1 && sellSignal) || (state == -1 && buySignal) {
		side := closeOrderSide(state)
		placeOrder(symbol, side)

		profit := (c.Close - entryPrice) * size
		if state == -1 {
			profit = (entryPrice - c.Close) * size
		}

		balance += size*c.Close + profit
		totalProfitLoss += profit

		if profit >= 0 {
			numOfWin++
		} else {
			numOfLose++
		}

		sendTelegramMessage(token,
			fmt.Sprintf("CLOSE [%s]\nPrice: %s\nPNL: %.2f\nBalance: %.2f\nTotal PNL: %.2f\nWin: %d Lose: %d",
				positionLabel(state), price, profit, balance, totalProfitLoss, numOfWin, numOfLose))

		state, entryPrice = 0, 0
	}
}

func placeOrder(symbol, side string) {
	endpoint := "https://fapi.binance.com/fapi/v1/order"
	ts := time.Now().UnixMilli()

	qty := strconv.FormatFloat(constant.QuantityMap[symbol], 'f',
		constant.SymbolPrecisionMap[symbol][1], 64)

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", "MARKET")
	params.Set("quantity", qty)
	params.Set("timestamp", strconv.FormatInt(ts, 10))

	signature := sign(params.Encode(), config.BinanceApiSecret)
	params.Set("signature", signature)

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBufferString(params.Encode()))
	req.Header.Set("X-MBX-APIKEY", config.BinanceApiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
}

// sign generates HMAC-SHA256 signature
func sign(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func calcRSI(closes []float64, length int) float64 {
	var gain, loss float64
	for i := len(closes) - length; i < len(closes); i++ {
		diff := closes[i] - closes[i-1]
		if diff > 0 {
			gain += diff
		} else {
			loss -= diff
		}
	}
	if loss == 0 {
		return 100
	}
	rs := gain / loss
	return 100 - (100 / (1 + rs))
}

func sma(data []float64, length int) float64 {
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

func parse(v any) float64 {
	f, _ := strconv.ParseFloat(v.(string), 64)
	return f
}
