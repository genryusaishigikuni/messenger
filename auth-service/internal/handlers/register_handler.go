package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/genryusaishigikuni/messenger/auth-service/internal/storage"
	"github.com/genryusaishigikuni/messenger/auth-service/pkg/utils"
)

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Info("Handling user registration request...")

		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error("Invalid request body: " + err.Error())
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		req.Username = strings.TrimSpace(req.Username)
		if req.Username == "" || req.Password == "" {
			utils.Error("Username or password is empty")
			http.Error(w, "Username and password required", http.StatusBadRequest)
			return
		}

		utils.Info("Checking if username already exists...")
		exists, err := storage.UserExists(db, req.Username)
		if err != nil {
			utils.Error("Error checking user existence: " + err.Error())
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		if exists {
			utils.Error("Username already taken: " + req.Username)
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}

		utils.Info("Hashing the user password...")
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.Error("Failed to hash password: " + err.Error())
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		utils.Info("Creating user in the database...")
		err = storage.CreateUser(db, req.Username, string(hashed))
		if err != nil {
			utils.Error("Failed to create user: " + err.Error())
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		utils.Info("User successfully registered: " + req.Username)
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(`{"message":"user registered"}`))
		if err != nil {
			utils.Error("Failed to write response: " + err.Error())
		}
	}
}
