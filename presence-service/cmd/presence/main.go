package main

import (
	"fmt"
	"net/http"

	"github.com/genryusaishigikuni/messenger/presence-service/internal/handlers"
	"github.com/genryusaishigikuni/messenger/presence-service/internal/memory"
	"github.com/genryusaishigikuni/messenger/presence-service/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	utils.Info("Loading configuration...")
	cfg := utils.LoadConfig()

	utils.Info("Initializing in-memory presence store...")
	store := memory.NewPresenceStore()
	utils.Info("Presence store initialized successfully.")

	utils.Info("Setting up routes...")
	r := mux.NewRouter()
	r.HandleFunc("/api/presence", handlers.GetPresenceHandler(store)).Methods("GET")
	utils.Info("Route set for GET /api/presence")
	r.HandleFunc("/api/presence/join", handlers.JoinHandler(store, cfg.AuthServiceURL)).Methods("POST")
	utils.Info("Route set for POST /api/presence/join")
	r.HandleFunc("/api/presence/leave", handlers.LeaveHandler(store, cfg.AuthServiceURL)).Methods("POST")
	utils.Info("Route set for POST /api/presence/leave")

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if req.Method == http.MethodOptions {
			utils.Info("OPTIONS request handled.")
			w.WriteHeader(http.StatusOK)
			return
		}

		utils.Info(fmt.Sprintf("Handling request: %s %s", req.Method, req.URL.Path))
		r.ServeHTTP(w, req)
	})

	utils.Info(fmt.Sprintf("Presence service starting on port %s...", cfg.ServerPort))
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		utils.Error(fmt.Sprintf("Server failed to start: %v", err))
	}
}
