package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BinanceApiKey    string
	BinanceApiSecret string
	BTCFuture4h      string
	BTCFuture1h      string
	ETHFuture4h      string
	ETHFuture1h      string
	ADAFuture4h      string
	BNBFuture4h      string
	SOLFuture4h      string
	TelegramChatId   string
)

// LoadConfig
func LoadConfig() {
	_ = godotenv.Load()

	BinanceApiKey = GetEnv("BINANCE_API_KEY")
	BinanceApiSecret = GetEnv("BINANCE_API_SECRET")
	BTCFuture4h = GetEnv("BTC_FUTURE_4h")
	BTCFuture1h = GetEnv("BTC_FUTURE_1h")
	ETHFuture4h = GetEnv("ETH_FUTURE_4h")
	ETHFuture1h = GetEnv("ETH_FUTURE_1h")
	ADAFuture4h = GetEnv("ADA_FUTURE_4h")
	BNBFuture4h = GetEnv("BNB_FUTURE_4h")
	SOLFuture4h = GetEnv("SOL_FUTURE_4h")
	TelegramChatId = GetEnv("TELEGRAM_CHAT_ID")
}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("%s environment variable not set", key)
	}
	return value
}
