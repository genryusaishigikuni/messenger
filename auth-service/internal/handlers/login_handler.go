package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/genryusaishigikuni/messenger/auth-service/internal/jwt"
	"github.com/genryusaishigikuni/messenger/auth-service/internal/storage"
	"github.com/genryusaishigikuni/messenger/auth-service/pkg/utils"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(db *sql.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Info("Handling login request...")

		var req loginRequest
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

		utils.Info("Fetching user from database...")
		user, err := storage.GetUserByUsername(db, req.Username)
		if err != nil {
			utils.Error("User not found or database error: " + err.Error())
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		utils.Info("Validating password...")
		if bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)) != nil {
			utils.Error("Password validation failed")
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		utils.Info("Generating JWT token...")
		token, err := jwt.GenerateToken(jwtSecret, user.ID, user.Username)
		if err != nil {
			utils.Error("Failed to generate JWT token: " + err.Error())
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		utils.Info("Login successful, responding with token")
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"token":"` + token + `"}`))
		if err != nil {
			utils.Error("Failed to write response: " + err.Error())
		}
	}
}
