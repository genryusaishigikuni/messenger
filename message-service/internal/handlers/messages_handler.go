package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/genryusaishigikuni/messenger/message-service/internal/storage"
)

// GetMessagesHandler GET /api/messages/history?channel=<id>
func GetMessagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channelIDStr := r.URL.Query().Get("channel")
		if channelIDStr == "" {
			http.Error(w, "channel query param required", http.StatusBadRequest)
			return
		}

		channelID, err := strconv.Atoi(channelIDStr)
		if err != nil {
			http.Error(w, "invalid channel id", http.StatusBadRequest)
			return
		}

		messages, err := storage.GetMessagesByChannel(db, channelID)
		if err != nil {
			http.Error(w, "could not retrieve messages", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(messages)
		if err != nil {
			return
		}
	}
}

// POST /api/messages { "channel_id": X, "content": "Hello" }
type createMessageRequest struct {
	ChannelID int    `json:"channel_id"`
	Content   string `json:"content"`
}

func CreateMessageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user_id from token (in a real scenario, validate token!)
		userID, err := extractUserIDFromToken(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req createMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if req.ChannelID < 1 || req.Content == "" {
			http.Error(w, "channel_id and content are required", http.StatusBadRequest)
			return
		}

		msg, err := storage.CreateMessage(db, req.ChannelID, userID, req.Content)
		if err != nil {
			http.Error(w, "could not create message", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(msg)
		if err != nil {
			return
		}
	}
}

// Placeholder function: In a real-world scenario, validate the JWT by calling Auth Service
type authValidateResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Valid    bool   `json:"valid"`
}

// extractUserIDFromToken now calls the Auth Service to validate the JWT.
func extractUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("no authorization header provided")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, errors.New("invalid authorization header format")
	}
	token := parts[1]

	// Validate token via Auth Service
	userID, err := validateTokenWithAuthService(token)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// validateTokenWithAuthService makes a GET request to Auth Service's /api/auth/validate endpoint
// with the bearer token and parses the response.
func validateTokenWithAuthService(token string) (int, error) {
	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "http://localhost:8082"
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", authURL+"/api/auth/validate", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call auth service: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("token validation failed with status %d", resp.StatusCode)
	}

	var validateResp authValidateResponse
	if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
		return 0, fmt.Errorf("failed to parse auth service response: %w", err)
	}

	if !validateResp.Valid {
		return 0, errors.New("invalid token according to auth service")
	}

	return validateResp.UserID, nil
}
