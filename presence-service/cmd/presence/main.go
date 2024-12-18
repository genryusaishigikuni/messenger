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

	log.Printf("Presence service running on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
