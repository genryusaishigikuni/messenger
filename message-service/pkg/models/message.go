package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	ChannelID int       `json:"channel_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
