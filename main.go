package main

import (
	"bot-1/bot"
	"bot-1/config"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// load config
	config.LoadConfig()

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter trading pair (e.g., btcusdt): ")
	symbolInput, _ := reader.ReadString('\n')
	symbol := strings.TrimSpace(symbolInput)

	fmt.Printf("Enter interval for %s (e.g., 1m, 5m, 15m, 1h): ", symbol)
	intervalInput, _ := reader.ReadString('\n')
	interval := strings.TrimSpace(intervalInput)

	fmt.Print("Enter stop loss percentage (e.g., 1.5): ")
	slInput, _ := reader.ReadString('\n')
	slStr := strings.TrimSpace(slInput)
	stopLossPercent, err := strconv.ParseFloat(slStr, 64)
	if err != nil {
		fmt.Println("Invalid stop loss input. Using default 1.5%.")
		stopLossPercent = 1.5
	}

	var token string
	var chatId string

	if interval == "1m" {
		token = config.TelegramTokenBTC1m
		chatId = config.TelegramChatId
	} else if interval == "5m" {
		token = config.TelegramToken
		chatId = config.TelegramChatId
	} else if interval == "15m" {

	} else if interval == "30m" {

	} else if interval == "1h" {

	} else if interval == "4h" {

	} else if interval == "1d" {

	}

	bot.Bot(symbol, interval, token, chatId, stopLossPercent)
}
