package utils

import "os"

type Config struct {
	ServerPort     string
	AuthServiceURL string
}

func LoadConfig() Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8083"
	}

	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "http://localhost:8082"
	}

	return Config{
		ServerPort:     port,
		AuthServiceURL: authURL,
	}
}
