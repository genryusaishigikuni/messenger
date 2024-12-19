package messageclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/genryusaishigikuni/messenger/gateway-service/pkg/models"
	"github.com/genryusaishigikuni/messenger/gateway-service/pkg/utils"
)

type createMessageRequest struct {
	ChannelID int    `json:"channel_id"`
	Content   string `json:"content"`
}

func CreateMessage(messageURL, token string, userID, channelID int, content string) (models.Message, error) {
	utils.Info("Preparing to create a message")
	_ = userID
	// Determine the message service URL
	if messageURL == "" {
		envURL := os.Getenv("MESSAGE_SERVICE_URL")
		if envURL == "" {
			envURL = "http://localhost:8081"
		}
		messageURL = envURL
	}
	utils.Info("Message service URL: " + messageURL)

	// Prepare the request data
	reqData := createMessageRequest{
		ChannelID: channelID,
		Content:   content,
	}
	jsonBytes, err := json.Marshal(reqData)
	if err != nil {
		utils.Error("Failed to marshal request data: " + err.Error())
		return models.Message{}, err
	}

	// Create an HTTP client and request
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", messageURL+"/api/messages", bytes.NewReader(jsonBytes))
	if err != nil {
		utils.Error("Failed to create HTTP request: " + err.Error())
		return models.Message{}, err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	utils.Info("HTTP request prepared with Authorization header")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		utils.Error("Failed to send request to message service: " + err.Error())
		return models.Message{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Error("Failed to close response body: " + err.Error())
		}
	}(resp.Body)

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		utils.Error(fmt.Sprintf("Message service returned status %d", resp.StatusCode))
		return models.Message{}, fmt.Errorf("failed to create message: status %d", resp.StatusCode)
	}
	utils.Info("Message service returned a successful response")

	// Decode the response body
	var msg models.Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		utils.Error("Failed to decode response body: " + err.Error())
		return models.Message{}, err
	}

	utils.Info("Message created successfully: " + fmt.Sprintf("ID=%d, ChannelID=%d, Content=%s", msg.ID, msg.ChannelID, msg.Content))
	return msg, nil
}
