package handlers

import (
	"encoding/json"
	"github.com/genryusaishigikuni/messenger/presence-service/internal/memory"
	"github.com/genryusaishigikuni/messenger/presence-service/pkg/utils"
	"net/http"
)

func GetPresenceHandler(store *memory.PresenceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Info("Handling get presence request...")

		// Retrieve all online users
		utils.Info("Fetching all online users from presence store")
		presences := store.GetAll()

		// Encode and send the response
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"online_users": presences,
		})
		if err != nil {
			utils.Error("Failed to encode response: " + err.Error())
			http.Error(w, "failed to retrieve online users", http.StatusInternalServerError)
			return
		}

		utils.Info("Presence information successfully sent to client")
	}
}
