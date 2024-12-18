package broadcaster

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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
	gatewayURL := os.Getenv("GATEWAY_SERVICE_URL")
	if gatewayURL == "" {
		gatewayURL = "http://localhost:8080"
	}

	ev := PresenceEvent{
		Event:     event,
		UserID:    userID,
		ChannelID: channelID,
	}

	data, err := json.Marshal(ev)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal presence event: %v", err)
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(gatewayURL+"/api/presence/event", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("[ERROR] Failed to send presence event to gateway: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("[ERROR] Failed to close presence event body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Gateway returned status %d for presence event", resp.StatusCode)
	} else {
		log.Println("[INFO] Successfully broadCasted presence event to gateway")
	}
}
