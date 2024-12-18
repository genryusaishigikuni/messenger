package handlers

import (
	"encoding/json"
	"github.com/genryusaishigikuni/messenger/presence-service/internal/broadcaster"
	"net/http"

	"github.com/genryusaishigikuni/messenger/presence-service/internal/memory"
)

type leaveRequest struct {
	// You can require only the token for user_id, but we have the token to identify the user anyway.
	// Optionally, user can specify if they leave a specific channel or just go offline completely.
	// Here we assume leaving means going completely offline.
}

func LeaveHandler(store *memory.PresenceStore, authServiceURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := extractUserIDFromToken(r, authServiceURL)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req leaveRequest
		_ = json.NewDecoder(r.Body).Decode(&req) // not needed, no fields currently

		// Get user presence before removing
		presence := store.GetPresence(userID)

		// Remove user from presence store
		store.SetOffline(userID)

		// Broadcast the event to the gateway or other listeners
		// If we knew the channel, use presence.ChannelID, else 0
		var channelID int
		if presence != nil {
			channelID = presence.ChannelID
		}
		broadcaster.BroadcastEvent("user_left", userID, channelID)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"message":"user left"}`))
		if err != nil {
			return
		}
	}
}
