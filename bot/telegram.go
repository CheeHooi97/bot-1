package bot

import (
	"bot-1/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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
	for {
		updates := fetchUpdates(token)
		for _, update := range updates {
			offset = update.UpdateID + 1 // Avoid reprocessing

			// Dispatch each update in its own goroutine
			go handleTelegramCommand(update, token, threadId)
		}
		time.Sleep(1 * time.Second) // Avoid hammering Telegram API
	}
}

func handleTelegramCommand(update Update, token string, threadId int64) {
	message := update.Message.Text

	switch message {
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
		log.Println("Unknown command:", message)
	}
}

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

var offset int = 0

func fetchUpdates(token string) []Update {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=5", token, offset)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching updates:", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result UpdatesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("JSON unmarshal error:", err)
		return nil
	}

	if !result.Ok {
		log.Println("getUpdates not OK:", string(body))
		return nil
	}

	return result.Result
}
