package storage

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/genryusaishigikuni/messenger/message-service/pkg/utils"
)

func RunMigrations(db *sql.DB, migrationsDir string) error {
	utils.Info(fmt.Sprintf("Starting to run migrations from directory: %s", migrationsDir))
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		utils.Error(fmt.Sprintf("Failed to read migrations directory: %v", err))
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			path := migrationsDir + "/" + entry.Name()
			utils.Info(fmt.Sprintf("Running migration: %s", path))
			content, err := os.ReadFile(path)
			if err != nil {
				utils.Error(fmt.Sprintf("Failed to read migration file %s: %v", path, err))
				return err
			}
			_, err = db.Exec(string(content))
			if err != nil {
				utils.Error(fmt.Sprintf("Failed to execute migration %s: %v", entry.Name(), err))
				return fmt.Errorf("failed to run migration %s: %v", entry.Name(), err)
			}
			utils.Info(fmt.Sprintf("Successfully applied migration: %s", path))
		}
	}

	utils.Info("All migrations applied successfully")
	return nil
}
