package handlers

import (
	"encoding/json"
	"net/http"
	_ "strconv"

	"github.com/genryusaishigikuni/messenger/presence-service/internal/memory"
)

func GetPresenceHandler(store *memory.PresenceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// No auth required for getting presence, or you can require auth if desired.
		presences := store.GetAll()
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"online_users": presences,
		})
		if err != nil {
			return
		}
	}
}
