package handlers

import (
	"encoding/json"
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

		// Possibly parse request body if needed
		var req leaveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// We don't really need info from body for now
		}

		store.SetOffline(userID)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"user left"}`))
	}
}
