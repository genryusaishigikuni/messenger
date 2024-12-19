package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/genryusaishigikuni/messenger/message-service/internal/storage"
	"github.com/genryusaishigikuni/messenger/message-service/pkg/utils"
)

type createChannelRequest struct {
	Name string `json:"name"`
}

func GetChannelsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Info("Received request to get channels")

		channels, err := storage.GetChannels(db)
		if err != nil {
			utils.Error("Failed to retrieve channels: " + err.Error())
			http.Error(w, "could not retrieve channels", http.StatusInternalServerError)
			return
		}

		utils.Info("Channels retrieved successfully")
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"channels": channels,
		})
		if err != nil {
			utils.Error("Failed to encode channels response: " + err.Error())
			return
		}
		utils.Info("Channels response sent successfully")
	}
}

func CreateChannelHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Info("Received request to create a channel")

		var req createChannelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error("Invalid create channel request: " + err.Error())
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			utils.Error("Channel name is missing in the request")
			http.Error(w, "channel name required", http.StatusBadRequest)
			return
		}

		utils.Info("Creating channel with name: " + req.Name)
		channel, err := storage.CreateChannel(db, req.Name)
		if err != nil {
			utils.Error("Failed to create channel: " + err.Error())
			http.Error(w, "could not create channel", http.StatusInternalServerError)
			return
		}

		utils.Info("Channel created successfully: " + channel.Name)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(channel)
		if err != nil {
			utils.Error("Failed to encode create channel response: " + err.Error())
			return
		}
		utils.Info("Create channel response sent successfully")
	}
}
