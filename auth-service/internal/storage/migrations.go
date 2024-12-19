package storage

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/genryusaishigikuni/messenger/auth-service/pkg/utils"
)

func RunMigrations(db *sql.DB, migrationsDir string) error {
	utils.Info("Starting database migrations...")

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		utils.Error("Failed to read migrations directory: " + err.Error())
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			path := migrationsDir + "/" + entry.Name()
			utils.Info("Running migration: " + entry.Name())

			content, err := os.ReadFile(path)
			if err != nil {
				utils.Error("Failed to read migration file " + entry.Name() + ": " + err.Error())
				return err
			}

			_, err = db.Exec(string(content))
			if err != nil {
				utils.Error(fmt.Sprintf("Failed to run migration %s: %v", entry.Name(), err))
				return fmt.Errorf("failed to run migration %s: %v", entry.Name(), err)
			}

			utils.Info("Successfully ran migration: " + entry.Name())
		}
	}

	utils.Info("All migrations completed successfully.")
	return nil
}
