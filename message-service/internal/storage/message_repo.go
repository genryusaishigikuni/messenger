package storage

import (
	"database/sql"
	"log"
	"time"

	"github.com/genryusaishigikuni/messenger/message-service/pkg/models"
)

func CreateMessage(db *sql.DB, channelID, userID int, content string) (*models.Message, error) {
	res, err := db.Exec("INSERT INTO messages (channel_id, user_id, content) VALUES (?, ?, ?)", channelID, userID, content)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Message{
		ID:        int(id),
		ChannelID: channelID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

func GetMessagesByChannel(db *sql.DB, channelID int) ([]models.Message, error) {
	rows, err := db.Query("SELECT id, channel_id, user_id, content, created_at FROM messages WHERE channel_id = ? ORDER BY created_at ASC", channelID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		err := rows.Scan(&m.ID, &m.ChannelID, &m.UserID, &m.Content, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
