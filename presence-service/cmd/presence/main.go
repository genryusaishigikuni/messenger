package main

import (
	"log"
	"net/http"

	"github.com/genryusaishigikuni/messenger/presence-service/internal/handlers"
	"github.com/genryusaishigikuni/messenger/presence-service/internal/memory"
	"github.com/genryusaishigikuni/messenger/presence-service/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	cfg := utils.LoadConfig()

	store := memory.NewPresenceStore()

	r := mux.NewRouter()
	r.HandleFunc("/api/presence", handlers.GetPresenceHandler(store)).Methods("GET")
	r.HandleFunc("/api/presence/join", handlers.JoinHandler(store, cfg.AuthServiceURL)).Methods("POST")
	r.HandleFunc("/api/presence/leave", handlers.LeaveHandler(store, cfg.AuthServiceURL)).Methods("POST")

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		r.ServeHTTP(w, req)
	})

	log.Printf("Presence service running on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
