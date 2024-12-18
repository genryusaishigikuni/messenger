package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type authValidateResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Valid    bool   `json:"valid"`
}

func extractUserIDFromToken(r *http.Request, authServiceURL string) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("no authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, errors.New("invalid authorization header")
	}
	token := parts[1]

	return validateTokenWithAuthService(token, authServiceURL)
}

func validateTokenWithAuthService(token, authURL string) (int, error) {
	if authURL == "" {
		// fallback to environment if not provided
		envURL := os.Getenv("AUTH_SERVICE_URL")
		if envURL == "" {
			envURL = "http://localhost:8082"
		}
		authURL = envURL
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
			log.Fatal(err)
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
