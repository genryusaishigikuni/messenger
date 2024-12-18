package handlers

import (
	_ "database/sql"
	"encoding/json"
	_ "errors"
	_ "fmt"
	"github.com/genryusaishigikuni/messenger/presence-service/internal/broadcaster"
	"net/http"
	_ "os"
	_ "strings"
	_ "time"

	"github.com/genryusaishigikuni/messenger/presence-service/internal/memory"
)

type joinRequest struct {
	ChannelID int `json:"channel_id"`
}

func JoinHandler(store *memory.PresenceStore, authServiceURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := extractUserIDFromToken(r, authServiceURL)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req joinRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if req.ChannelID < 1 {
			http.Error(w, "invalid channel_id", http.StatusBadRequest)
			return
		}

		store.SetOnline(userID, req.ChannelID)

		// Broadcast the event to the gateway or other listeners
		broadcaster.BroadcastEvent("user_joined", userID, req.ChannelID)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"message":"user joined"}`))
		if err != nil {
			return
		}
	}
}
