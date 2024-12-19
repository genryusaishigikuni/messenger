package main

import (
	"database/sql"
	"net/http"

	"github.com/genryusaishigikuni/messenger/auth-service/internal/handlers"
	"github.com/genryusaishigikuni/messenger/auth-service/internal/storage"
	"github.com/genryusaishigikuni/messenger/auth-service/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	// Load config
	utils.Info("Loading configuration...")
	cfg := utils.LoadConfig()

	// Initialize DB
	utils.Info("Initializing database...")
	db, err := storage.InitDB(cfg.DatabasePath)
	if err != nil {
		utils.Error("Failed to initialize database: " + err.Error())
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			utils.Error("Failed to close database connection: " + err.Error())
		}
	}(db)

	// Run migrations (simple approach)
	utils.Info("Running database migrations...")
	if err := storage.RunMigrations(db, "./migrations"); err != nil {
		utils.Error("Failed to run migrations: " + err.Error())
		panic(err)
	}

	// Prepare router
	utils.Info("Setting up routes...")
	r := mux.NewRouter()

	// Handlers
	r.HandleFunc("/api/auth/register", handlers.RegisterHandler(db)).Methods("POST")
	r.HandleFunc("/api/auth/login", handlers.LoginHandler(db, cfg.JWTSecret)).Methods("POST")
	r.HandleFunc("/api/auth/validate", handlers.ValidateHandler(cfg.JWTSecret)).Methods("GET")

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

	utils.Info("Starting Auth service on port " + cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		utils.Error("Server failed: " + err.Error())
	}
}
