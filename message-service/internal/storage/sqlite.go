package storage

import (
	"database/sql"
	"fmt"

	"github.com/genryusaishigikuni/messenger/message-service/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) (*sql.DB, error) {
	utils.Info(fmt.Sprintf("Initializing database connection to: %s", path))
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		utils.Error(fmt.Sprintf("Failed to open database at %s: %v", path, err))
		return nil, err
	}
	utils.Info("Database connection initialized successfully")
	return db, nil
}
