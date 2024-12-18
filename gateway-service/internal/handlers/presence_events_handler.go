package handlers

import (
	"encoding/json"
	"net/http"
)

type presenceEventRequest struct {
	Event     string `json:"event"` // "user_joined" or "user_left"
	UserID    int    `json:"user_id"`
	ChannelID int    `json:"channel_id"`
}

func PresenceEventHandler(manager *ConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ev presenceEventRequest
		if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Broadcast presence event to all connected clients in that channel
		manager.BroadcastPresenceEvent(ev.Event, ev.UserID, ev.ChannelID)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message":"received"}`))
		if err != nil {
			return
		}
	}
}
