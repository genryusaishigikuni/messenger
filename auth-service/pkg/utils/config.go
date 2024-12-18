package utils

import (
	"os"
)

type Config struct {
	DatabasePath string
	JWTSecret    string
	ServerPort   string
}

func LoadConfig() Config {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./auth.db"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}
	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	return Config{
		DatabasePath: dbPath,
		JWTSecret:    jwtSecret,
		ServerPort:   port,
	}
}
