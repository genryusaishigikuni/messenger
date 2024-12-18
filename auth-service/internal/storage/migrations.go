package storage

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func RunMigrations(db *sql.DB, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// Ensure it's a file (not a directory) and ends with .sql
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			path := migrationsDir + "/" + entry.Name()
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = db.Exec(string(content))
			if err != nil {
				return fmt.Errorf("failed to run migration %s: %v", entry.Name(), err)
			}
		}
	}
	return nil
}
