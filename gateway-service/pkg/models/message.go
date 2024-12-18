package models

import (
	"encoding/json"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	ChannelID int       `json:"channel_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// MarshalJSON is just the default, but let's just rely on the default marshaller.
func (m Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        int       `json:"id"`
		ChannelID int       `json:"channel_id"`
		UserID    int       `json:"user_id"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		UserID:    m.UserID,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
	})
}

// Generic JSON marshaller for map[string]interface{} if needed

func MarshalJSONGeneric(data map[string]interface{}) ([]byte, error) {
	return json.Marshal(data)
}
