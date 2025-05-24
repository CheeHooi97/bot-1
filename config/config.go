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
	// TelegramTokenBTC5m  string
	TelegramTokenBTC15m string
	// TelegramTokenBTC30m string
	TelegramTokenBTC1h string
	TelegramTokenBTC4h string
	// TelegramTokenBTC1d  string
	TelegramTokenETH15m string
	TelegramTokenETH4h  string
	TelegramChatId      string
	TelegramTestToken   string
	TelegramTestToken2  string
)

// LoadConfig
func LoadConfig() {
	_ = godotenv.Load()

	BinanceApiKey = GetEnv("BINANCE_API_KEY")
	BinanceApiSecret = GetEnv("BINANCE_API_SECRET")
	TelegramToken = GetEnv("TELEGRAM_TOKEN")
	TelegramTokenBTC1m = GetEnv("TELEGRAM_TOKEN_BTC_1m")
	// TelegramTokenBTC5m = GetEnv("TELEGRAM_TOKEN_BTC_5m")
	TelegramTokenBTC15m = GetEnv("TELEGRAM_TOKEN_BTC_15m")
	// TelegramTokenBTC30m = GetEnv("TELEGRAM_TOKEN_BTC_30m")
	TelegramTokenBTC1h = GetEnv("TELEGRAM_TOKEN_BTC_1h")
	TelegramTokenBTC4h = GetEnv("TELEGRAM_TOKEN_BTC_4h")
	// TelegramTokenBTC1d = GetEnv("TELEGRAM_TOKEN_BTC_1d")
	TelegramTokenETH15m = GetEnv("TELEGRAM_TOKEN_ETH_15m")
	TelegramTokenETH4h = GetEnv("TELEGRAM_TOKEN_ETH_4h")
	TelegramChatId = GetEnv("TELEGRAM_CHAT_ID")
	TelegramTestToken = GetEnv("TELEGRAM_TEST_TOKEN")
	TelegramTestToken2 = GetEnv("TELEGRAM_TEST_TOKEN2")

}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("%s environment variable not set", key)
	}
	return value
}
