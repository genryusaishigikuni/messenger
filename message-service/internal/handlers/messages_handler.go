package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/genryusaishigikuni/messenger/message-service/internal/storage"
	"github.com/genryusaishigikuni/messenger/message-service/pkg/utils"
)

// GetMessagesHandler GET /api/messages/history?channel=<id>
func GetMessagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Info("Received request to get messages")
		channelIDStr := r.URL.Query().Get("channel")
		if channelIDStr == "" {
			utils.Error("Channel query parameter is missing")
			http.Error(w, "channel query param required", http.StatusBadRequest)
			return
		}

		channelID, err := strconv.Atoi(channelIDStr)
		if err != nil {
			utils.Error("Invalid channel ID")
			http.Error(w, "invalid channel id", http.StatusBadRequest)
			return
		}

		messages, err := storage.GetMessagesByChannel(db, channelID)
		if err != nil {
			utils.Error(fmt.Sprintf("Failed to retrieve messages: %v", err))
			http.Error(w, "could not retrieve messages", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(messages); err != nil {
			utils.Error("Failed to encode messages response")
			return
		}
		utils.Info("Successfully retrieved messages")
	}
}

// POST /api/messages { "channel_id": X, "content": "Hello" }
type createMessageRequest struct {
	ChannelID int    `json:"channel_id"`
	Content   string `json:"content"`
}

func CreateMessageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Info("Received request to create a message")
		userID, err := extractUserIDFromToken(r)
		if err != nil {
			utils.Error(fmt.Sprintf("Unauthorized request: %v", err))
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req createMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error("Invalid request body")
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if req.ChannelID < 1 || req.Content == "" {
			utils.Error("Channel ID or content is missing")
			http.Error(w, "channel_id and content are required", http.StatusBadRequest)
			return
		}

		msg, err := storage.CreateMessage(db, req.ChannelID, userID, req.Content)
		if err != nil {
			utils.Error(fmt.Sprintf("Failed to create message: %v", err))
			http.Error(w, "could not create message", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(msg); err != nil {
			utils.Error("Failed to encode create message response")
			return
		}
		utils.Info("Message created successfully")
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
	utils.Info("Extracting user ID from token")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.Error("No authorization header provided")
		return 0, errors.New("no authorization header provided")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.Error("Invalid authorization header format")
		return 0, errors.New("invalid authorization header format")
	}
	token := parts[1]

	userID, err := validateTokenWithAuthService(token)
	if err != nil {
		utils.Error(fmt.Sprintf("Token validation failed: %v", err))
		return 0, err
	}

	utils.Info("Token validated successfully")
	return userID, nil
}

// validateTokenWithAuthService makes a GET request to Auth Service's /api/auth/validate endpoint
// with the bearer token and parses the response.
func validateTokenWithAuthService(token string) (int, error) {
	utils.Info("Validating token with Auth Service")
	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "http://localhost:8082"
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", authURL+"/api/auth/validate", nil)
	if err != nil {
		utils.Error(fmt.Sprintf("Failed to create request: %v", err))
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		utils.Error(fmt.Sprintf("Failed to call auth service: %v", err))
		return 0, fmt.Errorf("failed to call auth service: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if chError := Body.Close(); chError != nil {
			utils.Error("Failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		utils.Error(fmt.Sprintf("Token validation failed with status %d", resp.StatusCode))
		return 0, fmt.Errorf("token validation failed with status %d", resp.StatusCode)
	}

	var validateResp authValidateResponse
	if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
		utils.Error(fmt.Sprintf("Failed to parse auth service response: %v", err))
		return 0, fmt.Errorf("failed to parse auth service response: %w", err)
	}

	if !validateResp.Valid {
		utils.Error("Invalid token according to auth service")
		return 0, errors.New("invalid token according to auth service")
	}

	utils.Info("Token validated successfully by Auth Service")
	return validateResp.UserID, nil
}
