package main

import (
	"bot-1/config"
	"bot-1/handler"
	"fmt"
)

func main() {
	config.LoadConfig()

	h := handler.NewDCAHandler()

	fmt.Println("Starting manual trading CLI...")

	if err := h.StartDCA(); err != nil {
		fmt.Println("Error:", err)
	}

	select {}
}
