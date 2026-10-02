package notification

import (
	"encoding/json"
	"fmt"
	"log"
	"bytes"
	"io"
	"net/http"
)

type TelegramNotification struct {
	ChatID string
	Token string
	Message string
}

func (t TelegramNotification) Send(){
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage",t.Token)

	body,err := json.Marshal(map[string]string{
		"chat_id": t.ChatID,
		"text": t.Message,
	})
	if err != nil {
		log.Printf("Error Marshaling the request %v", err)
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error Sending message to telegram api: %v", err)
	}
	defer response.Body.Close()
	bodyResponse,err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response from telegram api: %v", err)
	}
	log.Printf("Message: %s \n set to telegram api", t.Message)
	log.Printf("Telegram API response: %s", string(bodyResponse))

}



