package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func sendTelegramMessage(token, message string, threadId int64) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	data := map[string]any{
		"chat_id":           -1002415528260,
		"message_thread_id": threadId,
		"text":              message,
	}

	payload, _ := json.Marshal(data)

	resp, err := http.Post(
		apiURL,
		"application/x-www-form-urlencoded",
		bytes.NewBuffer(payload),
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
