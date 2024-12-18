package utils

import "os"

type Config struct {
	AuthServiceURL    string
	MessageServiceURL string
	ServerPort        string
}

func LoadConfig() Config {
	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "http://localhost:8082"
	}

	msgURL := os.Getenv("MESSAGE_SERVICE_URL")
	if msgURL == "" {
		msgURL = "http://localhost:8081"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		AuthServiceURL:    authURL,
		MessageServiceURL: msgURL,
		ServerPort:        port,
	}
}
