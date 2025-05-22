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

	fmt.Print("Enter number of symbols to monitor: ")
	numInput, _ := reader.ReadString('\n')
	numInput = strings.TrimSpace(numInput)
	numSymbols, err := strconv.Atoi(numInput)
	if err != nil || numSymbols <= 0 {
		fmt.Println("Invalid number")
		return
	}

	for i := 0; i < numSymbols; i++ {
		fmt.Printf("Enter trading pair #%d (e.g., btcusdt): ", i+1)
		symbolInput, _ := reader.ReadString('\n')
		symbol := strings.TrimSpace(symbolInput)

		fmt.Printf("Enter interval for %s (e.g., 1m, 5m, 15m, 1h): ", symbol)
		intervalInput, _ := reader.ReadString('\n')
		interval := strings.TrimSpace(intervalInput)

		go bot.StartBot(symbol, interval) // run each symbol in a goroutine
	}

	// Keep the main function alive
	fmt.Println("Bots started. Press Ctrl+C to stop.")
	select {} // block forever
}
