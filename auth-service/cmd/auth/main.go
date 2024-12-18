package main

import (
	"database/sql"
	"log"
	"net/http"
	_ "os"

	"github.com/genryusaishigikuni/messenger/auth-service/internal/handlers"
	"github.com/genryusaishigikuni/messenger/auth-service/internal/storage"
	"github.com/genryusaishigikuni/messenger/auth-service/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	// Load config
	cfg := utils.LoadConfig()

	// Initialize DB
	db, err := storage.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	// Run migrations (simple approach)
	if err := storage.RunMigrations(db, "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Prepare router
	r := mux.NewRouter()

	// Handlers
	r.HandleFunc("/api/auth/register", handlers.RegisterHandler(db)).Methods("POST")
	r.HandleFunc("/api/auth/login", handlers.LoginHandler(db, cfg.JWTSecret)).Methods("POST")
	r.HandleFunc("/api/auth/validate", handlers.ValidateHandler(cfg.JWTSecret)).Methods("GET")

	log.Printf("Auth service running on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
