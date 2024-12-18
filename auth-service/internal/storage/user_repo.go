package storage

import (
	"database/sql"
	"errors"

	"github.com/genryusaishigikuni/messenger/auth-service/pkg/models"
)

func CreateUser(db *sql.DB, username, hashedPassword string) error {
	_, err := db.Exec("INSERT INTO users (username, hashed_password) VALUES (?, ?)", username, hashedPassword)
	return err
}

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	row := db.QueryRow("SELECT id, username, hashed_password, created_at FROM users WHERE username = ?", username)
	u := &models.User{}
	err := row.Scan(&u.ID, &u.Username, &u.HashedPassword, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}
	return u, nil
}

func UserExists(db *sql.DB, username string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
