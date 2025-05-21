package main

import (
	"bot-1/bot"
	"bot-1/config"
	"log"
)

func main() {
	// load config
	config.LoadConfig()

	log.Println("Test:")
	log.Println("Key:", config.BinanceApiKey)
	log.Println("Secret:", config.BinanceApiSecret)

	bot.BTC()
	// bot.ETH()
	// select {}
}
