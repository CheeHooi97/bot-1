package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OrderBook struct {
	LastUpdateID int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"` // [price, quantity]
	Asks         [][]string `json:"asks"`
}

func TestFuture(symbol string) {
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/depth?symbol=%s", symbol)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch depth: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: status code %d", resp.StatusCode)
	}

	var orderBook OrderBook
	if err := json.NewDecoder(resp.Body).Decode(&orderBook); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	log.Println("Order Book:")
	log.Println("Bids:")
	for _, bid := range orderBook.Bids {
		log.Printf("Price: %s, Quantity: %s\n", bid[0], bid[1])
	}

	log.Println("Asks:")
	for _, ask := range orderBook.Asks {
		log.Printf("Price: %s, Quantity: %s\n", ask[0], ask[1])
	}
}
