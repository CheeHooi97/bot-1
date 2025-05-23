package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BinanceApiKey      string
	BinanceApiSecret   string
	TelegramToken      string
	TelegramTokenBTC1m string
	TelegramChatId     string
)

// LoadConfig
func LoadConfig() {
	_ = godotenv.Load()

	BinanceApiKey = GetEnv("BINANCE_API_KEY")
	BinanceApiSecret = GetEnv("BINANCE_API_SECRET")
	TelegramToken = GetEnv("TELEGRAM_TOKEN")
	TelegramTokenBTC1m = GetEnv("TELEGRAM_TOKEN_BTC_1m")
	TelegramChatId = GetEnv("TELEGRAM_CHAT_ID")

}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("%s environment variable not set", key)
	}
	return value
}
