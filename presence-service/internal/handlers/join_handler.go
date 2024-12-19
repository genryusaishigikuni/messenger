package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/genryusaishigikuni/messenger/presence-service/internal/broadcaster"
	"github.com/genryusaishigikuni/messenger/presence-service/internal/memory"
	"github.com/genryusaishigikuni/messenger/presence-service/pkg/utils"
)

type joinRequest struct {
	ChannelID int `json:"channel_id"`
}

func JoinHandler(store *memory.PresenceStore, authServiceURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Info("Handling user join request...")

		// Extract user ID from token
		userID, err := extractUserIDFromToken(r, authServiceURL)
		if err != nil {
			utils.Error("Unauthorized access: " + err.Error())
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		utils.Info("Extracted user ID successfully: " + strconv.Itoa(userID))

		// Parse request body
		var req joinRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error("Failed to decode join request body: " + err.Error())
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		utils.Info("Decoded join request body successfully")

		// Validate channel ID
		if req.ChannelID < 1 {
			utils.Error("Invalid channel ID provided in join request")
			http.Error(w, "invalid channel_id", http.StatusBadRequest)
			return
		}
		utils.Info("Valid channel ID provided")

		// Update presence store
		store.SetOnline(userID, req.ChannelID)
		utils.Info("Updated presence store: User " + strconv.Itoa(userID) + " joined channel " + strconv.Itoa(req.ChannelID))

		// Broadcast the event
		utils.Info("Broadcasting user join event to gateway")
		broadcaster.BroadcastEvent("user_joined", userID, req.ChannelID)

		// Respond to the client
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"message":"user joined"}`))
		if err != nil {
			utils.Error("Failed to write response: " + err.Error())
			return
		}
		utils.Info("User join request handled successfully")
	}
}
