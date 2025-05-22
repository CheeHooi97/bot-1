package main

import (
	"bot-1/bot"
	"bot-1/config"
	"bufio"
	"fmt"
	"os"
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

	bot.Bot(symbol, interval)
}
