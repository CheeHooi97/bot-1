package bot

import (
	"bot-1/config"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func sendTelegramMessage(message string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramToken)
	data := url.Values{}
	data.Set("chat_id", config.TelegramChatId)
	data.Set("text", message)

	resp, err := http.Post(
		apiURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send message: %s\n", resp.Status)
	}
}
