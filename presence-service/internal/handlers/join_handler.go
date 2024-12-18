package handlers

import (
	_ "database/sql"
	"encoding/json"
	_ "errors"
	_ "fmt"
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

		// Mark user as online in the given channel
		store.SetOnline(userID, req.ChannelID)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"user joined"}`))
	}
}
