package handlers

import (
	"encoding/json"
	"github.com/genryusaishigikuni/messenger/presence-service/internal/broadcaster"
	"github.com/genryusaishigikuni/messenger/presence-service/internal/memory"
	"github.com/genryusaishigikuni/messenger/presence-service/pkg/utils"
	"net/http"
	"strconv"
)

type leaveRequest struct {
	// Placeholder struct for potential future use.
}

func LeaveHandler(store *memory.PresenceStore, authServiceURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Info("Handling user leave request...")

		// Extract user ID from token
		userID, err := extractUserIDFromToken(r, authServiceURL)
		if err != nil {
			utils.Error("Unauthorized access: " + err.Error())
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		utils.Info("Extracted user ID successfully: " + strconv.Itoa(userID))

		// Decode request body (even though it's unused)
		var req leaveRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.Error("Failed to decode leave request body: " + err.Error())
		} else {
			utils.Info("Decoded leave request body successfully (placeholder)")
		}

		// Get current presence information
		utils.Info("Retrieving presence information for user " + strconv.Itoa(userID))
		presence := store.GetPresence(userID)

		// Remove user from presence store
		utils.Info("Setting user offline in presence store: User " + strconv.Itoa(userID))
		store.SetOffline(userID)

		// Determine channel ID for broadcasting
		var channelID int
		if presence != nil {
			channelID = presence.ChannelID
			utils.Info("User was associated with channel ID: " + strconv.Itoa(channelID))
		} else {
			utils.Info("User was not associated with any channel")
		}

		// Broadcast the user left event
		utils.Info("Broadcasting user left event to gateway")
		broadcaster.BroadcastEvent("user_left", userID, channelID)

		// Respond to the client
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"message":"user left"}`))
		if err != nil {
			utils.Error("Failed to write response: " + err.Error())
			return
		}
		utils.Info("User leave request handled successfully")
	}
}
