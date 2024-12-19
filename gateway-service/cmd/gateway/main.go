package main

import (
	"log"
	"net/http"

	"github.com/genryusaishigikuni/messenger/gateway-service/internal/handlers"
	"github.com/genryusaishigikuni/messenger/gateway-service/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	utils.Info("Loading configuration")
	cfg := utils.LoadConfig()

	utils.Info("Initializing connection manager")
	manager := handlers.NewConnectionManager()

	utils.Info("Setting up router")
	r := mux.NewRouter()

	utils.Info("Registering WebSocket endpoint")
	// WebSocket endpoint
	r.HandleFunc("/ws", handlers.WebSocketHandler(manager, cfg.AuthServiceURL, cfg.MessageServiceURL))

	utils.Info("Registering Presence event endpoint")
	// Presence event endpoint (called by Presence Service)
	r.HandleFunc("/api/presence/event", handlers.PresenceEventHandler(manager)).Methods("POST")

	// Setting up middleware for CORS
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		utils.Info("Handling CORS for incoming request")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if req.Method == http.MethodOptions {
			utils.Info("Responding to OPTIONS preflight request")
			w.WriteHeader(http.StatusOK)
			return
		}

		r.ServeHTTP(w, req)
	})

	utils.Info("Starting Gateway service on port " + cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		utils.Error("Server failed: " + err.Error())
		log.Fatalf("Server failed: %v", err)
	}
}
