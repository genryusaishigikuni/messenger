package main

import (
	"log"
	"net/http"

	"github.com/genryusaishigikuni/messenger/gateway-service/internal/handlers"
	"github.com/genryusaishigikuni/messenger/gateway-service/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	cfg := utils.LoadConfig()

	manager := handlers.NewConnectionManager()

	r := mux.NewRouter()
	// WebSocket endpoint
	r.HandleFunc("/ws", handlers.WebSocketHandler(manager, cfg.AuthServiceURL, cfg.MessageServiceURL))

	// Presence event endpoint (called by Presence Service)
	r.HandleFunc("/api/presence/event", handlers.PresenceEventHandler(manager)).Methods("POST")

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

	log.Printf("Gateway service running on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
