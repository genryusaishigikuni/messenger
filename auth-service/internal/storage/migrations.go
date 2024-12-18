package storage

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func RunMigrations(db *sql.DB, migrationsDir string) error {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			path := migrationsDir + "/" + file.Name()
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = db.Exec(string(content))
			if err != nil {
				return fmt.Errorf("failed to run migration %s: %v", file.Name(), err)
			}
		}
	}
	return nil
}
