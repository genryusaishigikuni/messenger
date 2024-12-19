package main

import (
	"database/sql"
	"net/http"

	"github.com/genryusaishigikuni/messenger/message-service/internal/handlers"
	"github.com/genryusaishigikuni/messenger/message-service/internal/storage"
	"github.com/genryusaishigikuni/messenger/message-service/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	utils.Info("Starting Message Service...")

	// Load configuration
	utils.Info("Loading configuration...")
	cfg := utils.LoadConfig()
	utils.Info("Configuration loaded successfully")

	// Initialize database
	utils.Info("Initializing database connection...")
	db, err := storage.InitDB(cfg.DatabasePath)
	if err != nil {
		utils.Error("Failed to initialize database: " + err.Error())
		return
	}
	defer func(db *sql.DB) {
		utils.Info("Closing database connection...")
		err := db.Close()
		if err != nil {
			utils.Error("Failed to close database connection: " + err.Error())
		} else {
			utils.Info("Database connection closed successfully")
		}
	}(db)
	utils.Info("Database initialized successfully")

	// Run migrations
	utils.Info("Running database migrations...")
	if err := storage.RunMigrations(db, "./migrations"); err != nil {
		utils.Error("Failed to run migrations: " + err.Error())
		return
	}
	utils.Info("Database migrations completed successfully")

	// Setup router
	utils.Info("Setting up HTTP routes...")
	r := mux.NewRouter()

	// Channels endpoints
	r.HandleFunc("/api/channels", handlers.GetChannelsHandler(db)).Methods("GET")
	r.HandleFunc("/api/channels", handlers.CreateChannelHandler(db)).Methods("POST")

	// Messages endpoints
	r.HandleFunc("/api/messages/history", handlers.GetMessagesHandler(db)).Methods("GET")
	r.HandleFunc("/api/messages", handlers.CreateMessageHandler(db)).Methods("POST")

	// Add CORS support
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if req.Method == http.MethodOptions {
			utils.Info("CORS preflight request handled")
			w.WriteHeader(http.StatusOK)
			return
		}

		r.ServeHTTP(w, req)
	})

	// Start server
	utils.Info("Starting HTTP server on port " + cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		utils.Error("Server failed: " + err.Error())
	} else {
		utils.Info("Message service stopped gracefully")
	}
}
