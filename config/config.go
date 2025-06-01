package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BinanceApiKey    string
	BinanceApiSecret string
	BTC1m            string
	BTC5m            string
	BTC15m           string
	BTC1h            string
	BTC4h            string
	// BTC1d  string
	ETH1m  string
	ETH5m  string
	ETH15m string
	ETH1h  string
	ETH4h  string
	// ETH1d  string
	SOL1m          string
	SOL15m         string
	BNB1m          string
	BNB15m         string
	ADA1m          string
	ADA15m         string
	XRP1m          string
	XRP15m         string
	TelegramChatId string
)

// LoadConfig
func LoadConfig() {
	_ = godotenv.Load()

	BinanceApiKey = GetEnv("BINANCE_API_KEY")
	BinanceApiSecret = GetEnv("BINANCE_API_SECRET")
	BTC1m = GetEnv("BTC_1m")
	BTC5m = GetEnv("BTC_5m")
	BTC15m = GetEnv("BTC_15m")
	BTC1h = GetEnv("BTC_1h")
	BTC4h = GetEnv("BTC_4h")
	// BTC1d = GetEnv("BTC_1d")
	ETH1m = GetEnv("ETH_1m")
	ETH5m = GetEnv("ETH_5m")
	ETH15m = GetEnv("ETH_15m")
	ETH1h = GetEnv("ETH_1h")
	ETH4h = GetEnv("ETH_4h")
	SOL1m = GetEnv("SOL_1m")
	SOL15m = GetEnv("SOL_15m")
	BNB1m = GetEnv("BNB_1m")
	BNB15m = GetEnv("BNB_15m")
	ADA1m = GetEnv("ADA_1m")
	ADA15m = GetEnv("ADA_15m")
	XRP1m = GetEnv("XRP_1m")
	XRP15m = GetEnv("XRP_15m")
	TelegramChatId = GetEnv("TELEGRAM_CHAT_ID")
}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("%s environment variable not set", key)
	}
	return value
}
