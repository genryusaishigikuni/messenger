package storage

import (
	"database/sql"
	"log"
	"time"

	"github.com/genryusaishigikuni/messenger/message-service/pkg/models"
)

func CreateChannel(db *sql.DB, name string) (*models.Channel, error) {
	res, err := db.Exec("INSERT INTO channels (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Channel{
		ID:        int(id),
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
}

func GetChannels(db *sql.DB) ([]models.Channel, error) {
	rows, err := db.Query("SELECT id, name, created_at FROM channels ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var channels []models.Channel
	for rows.Next() {
		var c models.Channel
		err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		channels = append(channels, c)
	}
	return channels, nil
}
