package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/genryusaishigikuni/messenger/gateway-service/pkg/utils"
)

type presenceEventRequest struct {
	Event     string `json:"event"` // "user_joined" or "user_left"
	UserID    int    `json:"user_id"`
	ChannelID int    `json:"channel_id"`
}

func PresenceEventHandler(manager *ConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Info("Received presence event request")

		var ev presenceEventRequest
		if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
			utils.Error("Failed to decode presence event request: " + err.Error())
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		utils.Info("Broadcasting presence event: " + ev.Event + " for UserID=" + strconv.Itoa(ev.UserID) + " in ChannelID=" + strconv.Itoa(ev.ChannelID))

		// Broadcast presence event to all connected clients in that channel
		manager.BroadcastPresenceEvent(ev.Event, ev.UserID, ev.ChannelID)

		utils.Info("Presence event successfully broadCasted")

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message":"received"}`))
		if err != nil {
			utils.Error("Failed to send response: " + err.Error())
		}
	}
}
