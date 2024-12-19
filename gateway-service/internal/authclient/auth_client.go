package authclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/genryusaishigikuni/messenger/gateway-service/pkg/utils"
)

type validateResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Valid    bool   `json:"valid"`
}

// ValidateUserFromRequest validates a user from the Authorization header of a request.

// ValidateToken validates a token with the Auth Service.
func ValidateToken(token, authURL string) (int, error) {
	utils.Info("Starting token validation with Auth Service")
	if authURL == "" {
		utils.Info("Auth URL not provided, checking environment variable")
		envURL := os.Getenv("AUTH_SERVICE_URL")
		if envURL == "" {
			utils.Error("No AUTH_SERVICE_URL environment variable set, using default")
			envURL = "http://localhost:8082"
		}
		authURL = envURL
	}
	utils.Info("Using Auth Service URL: " + authURL)

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", authURL+"/api/auth/validate", nil)
	if err != nil {
		utils.Error("Failed to create request for token validation: " + err.Error())
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
			utils.Error("Failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		utils.Error(fmt.Sprintf("Token validation failed with status %d", resp.StatusCode))
		return 0, fmt.Errorf("token validation failed with status %d", resp.StatusCode)
	}

	var v validateResponse
	utils.Info("Parsing token validation response")
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		utils.Error("Failed to parse validate response: " + err.Error())
		return 0, fmt.Errorf("failed to parse validate response: %w", err)
	}

	if !v.Valid {
		utils.Error("Token is invalid according to Auth Service")
		return 0, errors.New("invalid token according to auth service")
	}

	utils.Info(fmt.Sprintf("Token validated successfully for user ID: %d", v.UserID))
	return v.UserID, nil
}
