package utils

import "os"

type Config struct {
	DatabasePath   string
	ServerPort     string
	AuthServiceURL string
}

func LoadConfig() Config {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./messages.db"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
	}

	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "http://localhost:8082"
	}

	return Config{
		DatabasePath:   dbPath,
		ServerPort:     port,
		AuthServiceURL: authURL,
	}
}
