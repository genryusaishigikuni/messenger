package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/genryusaishigikuni/messenger/presence-service/pkg/utils"
)

type authValidateResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Valid    bool   `json:"valid"`
}

func extractUserIDFromToken(r *http.Request, authServiceURL string) (int, error) {
	utils.Info("Extracting user ID from token...")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.Error("No authorization header provided")
		return 0, errors.New("no authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.Error("Invalid authorization header format")
		return 0, errors.New("invalid authorization header")
	}
	token := parts[1]
	utils.Info("Authorization header extracted successfully")

	return validateTokenWithAuthService(token, authServiceURL)
}

func validateTokenWithAuthService(token, authURL string) (int, error) {
	utils.Info("Validating token with Auth Service...")

	if authURL == "" {
		envURL := os.Getenv("AUTH_SERVICE_URL")
		if envURL == "" {
			envURL = "http://localhost:8082"
		}
		authURL = envURL
		utils.Info(fmt.Sprintf("Auth Service URL not provided. Using default: %s", authURL))
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", authURL+"/api/auth/validate", nil)
	if err != nil {
		utils.Error("Failed to create request: " + err.Error())
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	utils.Info("Sending token validation request to Auth Service")

	resp, err := client.Do(req)
	if err != nil {
		utils.Error("Failed to call Auth Service: " + err.Error())
		return 0, fmt.Errorf("failed to call auth service: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Error("Failed to close Auth Service response body: " + err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		utils.Error(fmt.Sprintf("Auth Service returned status %d", resp.StatusCode))
		return 0, fmt.Errorf("token validation failed with status %d", resp.StatusCode)
	}

	var validateResp authValidateResponse
	if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
		utils.Error("Failed to parse Auth Service response: " + err.Error())
		return 0, fmt.Errorf("failed to parse auth service response: %w", err)
	}

	if !validateResp.Valid {
		utils.Error("Token validation failed: invalid token")
		return 0, errors.New("invalid token according to auth service")
	}

	utils.Info(fmt.Sprintf("Token validated successfully for user ID: %d", validateResp.UserID))
	return validateResp.UserID, nil
}
