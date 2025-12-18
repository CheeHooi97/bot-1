package bot

import (
	"bot-1/config"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type DCABot struct {
	Symbol       string
	DropPercent  float64
	TotalUSDT    float64
	OneBuyUSDT   float64
	LastBuyPrice float64
	Started      bool
	Records      []DCARecord
}

type DCARecord struct {
	BuyNumber     int
	Price         float64
	USDTSpent     float64
	AmountBought  float64
	RemainingUSDT float64
	TotalHoldings float64
}
