package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/genryusaishigikuni/messenger/message-service/internal/storage"
)

type createChannelRequest struct {
	Name string `json:"name"`
}

func GetChannelsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channels, err := storage.GetChannels(db)
		if err != nil {
			http.Error(w, "could not retrieve channels", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"channels": channels,
		})
		if err != nil {
			return
		}
	}
}

func CreateChannelHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createChannelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			http.Error(w, "channel name required", http.StatusBadRequest)
			return
		}

		channel, err := storage.CreateChannel(db, req.Name)
		if err != nil {
			http.Error(w, "could not create channel", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(channel)
		if err != nil {
			return
		}
	}
}
