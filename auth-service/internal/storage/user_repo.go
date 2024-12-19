package storage

import (
	"database/sql"
	"errors"

	"github.com/genryusaishigikuni/messenger/auth-service/pkg/models"
	"github.com/genryusaishigikuni/messenger/auth-service/pkg/utils"
)

func CreateUser(db *sql.DB, username, hashedPassword string) error {
	utils.Info("Creating a new user: " + username)
	_, err := db.Exec("INSERT INTO users (username, hashed_password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		utils.Error("Failed to create user: " + err.Error())
		return err
	}
	utils.Info("User " + username + " created successfully.")
	return nil
}

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	utils.Info("Fetching user by username: " + username)
	row := db.QueryRow("SELECT id, username, hashed_password, created_at FROM users WHERE username = ?", username)
	u := &models.User{}
	err := row.Scan(&u.ID, &u.Username, &u.HashedPassword, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Error("User not found: " + username)
		return nil, errors.New("user not found")
	} else if err != nil {
		utils.Error("Failed to fetch user: " + err.Error())
		return nil, err
	}
	utils.Info("User " + username + " fetched successfully.")
	return u, nil
}

func UserExists(db *sql.DB, username string) (bool, error) {
	utils.Info("Checking if user exists: " + username)
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		utils.Error("Failed to check if user exists: " + err.Error())
		return false, err
	}
	if count > 0 {
		utils.Info("User exists: " + username)
	} else {
		utils.Info("User does not exist: " + username)
	}
	return count > 0, nil
}
