package storage

import (
	"database/sql"
	"github.com/genryusaishigikuni/messenger/auth-service/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) (*sql.DB, error) {
	utils.Info("Initializing database connection...")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		utils.Error("Failed to open database: " + err.Error())
		return nil, err
	}
	utils.Info("Database connection initialized successfully.")
	return db, nil
}
