package bot

import (
	"bot-1/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendTelegramMessage(token, message string, threadId int64) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	// Use a map for the JSON payload
	payload := map[string]any{
		"chat_id":           config.TelegramChatId, // Make sure it's correct
		"message_thread_id": threadId,              // Can be int64
		"text":              message,               // Must not be empty
	}

	// Marshal the map into JSON
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}

	// Send the request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send message: %s\nBody: %s\n", resp.Status, string(body))
	}
}

func listenTelegramCommands(token string, threadId int64) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}
		command := strings.ToLower(update.Message.Text)
		go handleTelegramCommand(command, token, threadId)
	}
}

func handleTelegramCommand(command string, token string, threadId int64) {
	switch command {
	case "/start":
		tradingEnabled = true
		msg := "Trading bot has been STARTED. Processing candles and placing orders."
		log.Println(msg)
		sendTelegramMessage(token, msg, threadId)
	case "/stop":
		tradingEnabled = false
		msg := "Trading bot has been STOPPED. Still processing candles, but will not trade."
		log.Println(msg)
		sendTelegramMessage(token, msg, threadId)
	default:
		log.Println("Unknown command:", command)
	}
}
