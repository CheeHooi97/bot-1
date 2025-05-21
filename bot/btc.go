package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

func BTC() {
	// client := binance.NewClient(config.BinanceApiKey, config.BinanceApiSecret)

	symbol := "btcusdt"
	interval := "15m"

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

			processCandle(candle, symbol)
		}
	}
}
