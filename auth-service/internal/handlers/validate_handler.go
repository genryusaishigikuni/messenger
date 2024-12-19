package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/genryusaishigikuni/messenger/auth-service/internal/jwt"
	"github.com/genryusaishigikuni/messenger/auth-service/pkg/utils"
)

func ValidateHandler(jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Info("Handling token validation request...")

		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Error("Authorization header is missing")
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		// Check Authorization header format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Error("Invalid authorization header format")
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Validate token
		tokenStr := parts[1]
		utils.Info("Validating token...")
		claims, err := jwt.ValidateToken(jwtSecret, tokenStr)
		if err != nil {
			utils.Error("Invalid token: " + err.Error())
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Send response
		utils.Info("Token validated successfully for user ID: " + strconv.Itoa(claims.UserID))
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"valid":    true,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			utils.Error("Failed to encode response: " + err.Error())
		}
	}
}
