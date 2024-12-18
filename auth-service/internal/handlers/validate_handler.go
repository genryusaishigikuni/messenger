package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/genryusaishigikuni/messenger/auth-service/internal/jwt"
)

func ValidateHandler(jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		claims, err := jwt.ValidateToken(jwtSecret, tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"valid":    true,
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}
	}
}
