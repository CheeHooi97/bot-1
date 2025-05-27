package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BinanceApiKey       string
	BinanceApiSecret    string
	TelegramTokenBTC1m  string
	TelegramTokenBTC5m  string
	TelegramTokenBTC15m string
	TelegramTokenBTC30m string
	TelegramTokenBTC1h  string
	TelegramTokenBTC4h  string
	// TelegramTokenBTC1d  string
	TelegramTokenETH1m  string
	TelegramTokenETH5m  string
	TelegramTokenETH15m string
	TelegramTokenETH30m string
	TelegramTokenETH1h  string
	TelegramTokenETH4h  string
	// TelegramTokenETH1d  string
	TelegramTokenSOL1m  string
	TelegramTokenSOL15m string
	TelegramTokenBNB1m  string
	TelegramTokenBNB15m string
	TelegramTokenADA1m  string
	TelegramTokenADA15m string
	TelegramTokenXRP1m  string
	TelegramTokenXRP15m string
	TelegramChatId      string
)

// LoadConfig
func LoadConfig() {
	_ = godotenv.Load()

	BinanceApiKey = GetEnv("BINANCE_API_KEY")
	BinanceApiSecret = GetEnv("BINANCE_API_SECRET")
	TelegramTokenBTC1m = GetEnv("TELEGRAM_TOKEN_BTC_1m")
	TelegramTokenBTC5m = GetEnv("TELEGRAM_TOKEN_BTC_5m")
	TelegramTokenBTC15m = GetEnv("TELEGRAM_TOKEN_BTC_15m")
	TelegramTokenBTC30m = GetEnv("TELEGRAM_TOKEN_BTC_30m")
	TelegramTokenBTC1h = GetEnv("TELEGRAM_TOKEN_BTC_1h")
	TelegramTokenBTC4h = GetEnv("TELEGRAM_TOKEN_BTC_4h")
	// TelegramTokenBTC1d = GetEnv("TELEGRAM_TOKEN_BTC_1d")
	TelegramTokenETH1m = GetEnv("TELEGRAM_TOKEN_ETH_1m")
	TelegramTokenETH5m = GetEnv("TELEGRAM_TOKEN_ETH_5m")
	TelegramTokenETH15m = GetEnv("TELEGRAM_TOKEN_ETH_15m")
	TelegramTokenETH30m = GetEnv("TELEGRAM_TOKEN_ETH_30m")
	TelegramTokenETH1h = GetEnv("TELEGRAM_TOKEN_ETH_1h")
	TelegramTokenETH4h = GetEnv("TELEGRAM_TOKEN_ETH_4h")
	TelegramTokenSOL1m = GetEnv("TELEGRAM_TOKEN_SOL_1m")
	TelegramTokenSOL15m = GetEnv("TELEGRAM_TOKEN_SOL_15m")
	TelegramTokenBNB1m = GetEnv("TELEGRAM_TOKEN_BNB_1m")
	TelegramTokenBNB15m = GetEnv("TELEGRAM_TOKEN_BNB_15m")
	TelegramTokenADA1m = GetEnv("TELEGRAM_TOKEN_ADA_1m")
	TelegramTokenADA15m = GetEnv("TELEGRAM_TOKEN_ADA_15m")
	TelegramTokenXRP1m = GetEnv("TELEGRAM_TOKEN_XRP_1m")
	TelegramTokenXRP15m = GetEnv("TELEGRAM_TOKEN_XRP_15m")
	TelegramChatId = GetEnv("TELEGRAM_CHAT_ID")
}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("%s environment variable not set", key)
	}
	return value
}
