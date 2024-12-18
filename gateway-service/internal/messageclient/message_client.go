package messageclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/genryusaishigikuni/messenger/gateway-service/pkg/models"
)

type createMessageRequest struct {
	ChannelID int    `json:"channel_id"`
	Content   string `json:"content"`
}

func CreateMessage(messageURL, token string, userID, channelID int, content string) (models.Message, error) {
	if messageURL == "" {
		envURL := os.Getenv("MESSAGE_SERVICE_URL")
		if envURL == "" {
			envURL = "http://localhost:8081"
		}
		messageURL = envURL
	}

	reqData := createMessageRequest{
		ChannelID: channelID,
		Content:   content,
	}
	jsonBytes, _ := json.Marshal(reqData)

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", messageURL+"/api/messages", bytes.NewReader(jsonBytes))
	if err != nil {
		return models.Message{}, err
	}

	// Include the user's JWT token so Message Service can validate user_id
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return models.Message{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return models.Message{}, fmt.Errorf("failed to create message: status %d", resp.StatusCode)
	}

	var msg models.Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return models.Message{}, err
	}

	return msg, nil
}
