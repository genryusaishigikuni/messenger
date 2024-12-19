package broadcaster

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/genryusaishigikuni/messenger/presence-service/pkg/utils"
)

type PresenceEvent struct {
	Event     string `json:"event"`
	UserID    int    `json:"user_id"`
	ChannelID int    `json:"channel_id,omitempty"`
}

// BroadcastEvent sends a presence event to the Gateway Service or another interested service.
// event can be "user_joined" or "user_left"
// userID is the ID of the user whose presence changed
// channelID is the channel the user is associated with (if applicable)
func BroadcastEvent(event string, userID, channelID int) {
	utils.Info("Preparing to broadcast presence event...")

	gatewayURL := os.Getenv("GATEWAY_SERVICE_URL")
	if gatewayURL == "" {
		gatewayURL = "http://localhost:8080"
		utils.Info("GATEWAY_SERVICE_URL not set. Using default: http://localhost:8080")
	}

	ev := PresenceEvent{
		Event:     event,
		UserID:    userID,
		ChannelID: channelID,
	}

	data, err := json.Marshal(ev)
	if err != nil {
		utils.Error("Failed to marshal presence event: " + err.Error())
		return
	}
	utils.Info("Presence event marshaled successfully.")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(gatewayURL+"/api/presence/event", "application/json", bytes.NewBuffer(data))
	if err != nil {
		utils.Error("Failed to send presence event to gateway: " + err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Error("Failed to close presence event response body: " + err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		utils.Error("Gateway returned status " + http.StatusText(resp.StatusCode) + " for presence event.")
	} else {
		utils.Info("Successfully broadCasted presence event to gateway.")
	}
}
