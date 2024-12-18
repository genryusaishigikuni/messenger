package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/genryusaishigikuni/messenger/message-service/internal/handlers"
	"github.com/genryusaishigikuni/messenger/message-service/internal/storage"
	"github.com/genryusaishigikuni/messenger/message-service/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	cfg := utils.LoadConfig()

	db, err := storage.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close DB: %v", err)
		}
	}(db)

	if err := storage.RunMigrations(db, "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	r := mux.NewRouter()

	// Channels endpoints
	r.HandleFunc("/api/channels", handlers.GetChannelsHandler(db)).Methods("GET")
	r.HandleFunc("/api/channels", handlers.CreateChannelHandler(db)).Methods("POST")

	// Messages endpoints
	r.HandleFunc("/api/messages/history", handlers.GetMessagesHandler(db)).Methods("GET")
	r.HandleFunc("/api/messages", handlers.CreateMessageHandler(db)).Methods("POST")

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

	log.Printf("Message service running on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
