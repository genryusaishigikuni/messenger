package handlers

import (
	"sync"

	"github.com/genryusaishigikuni/messenger/gateway-service/pkg/models"
	"github.com/gorilla/websocket"
)

type clientInfo struct {
	Conn      *websocket.Conn
	UserID    int
	Token     string
	ChannelID int
}

// ConnectionManager tracks channels and their connected clients
type ConnectionManager struct {
	mu       sync.RWMutex
	channels map[int]map[*websocket.Conn]*clientInfo // channel_id -> map[conn]*clientInfo
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		channels: make(map[int]map[*websocket.Conn]*clientInfo),
	}
}

func (m *ConnectionManager) RegisterClient(conn *websocket.Conn, userID, channelID int, token string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.channels[channelID]; !ok {
		m.channels[channelID] = make(map[*websocket.Conn]*clientInfo)
	}
	m.channels[channelID][conn] = &clientInfo{
		Conn:      conn,
		UserID:    userID,
		ChannelID: channelID,
		Token:     token,
	}
}

func (m *ConnectionManager) UnregisterClient(conn *websocket.Conn, channelID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if channel, ok := m.channels[channelID]; ok {
		delete(channel, conn)
		if len(channel) == 0 {
			delete(m.channels, channelID)
		}
	}
	err := conn.Close()
	if err != nil {
		return
	}
}

func (m *ConnectionManager) BroadcastToChannel(channelID int, msg models.Message) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	channel, ok := m.channels[channelID]
	if !ok {
		return
	}

	msgBytes, _ := msg.MarshalJSON()
	for _, client := range channel {
		err := client.Conn.WriteMessage(1, msgBytes)
		if err != nil {
			return
		}
	}
}

func (m *ConnectionManager) BroadcastPresenceEvent(eventType string, userID, channelID int) {
	event := map[string]interface{}{
		"event":      eventType,
		"user_id":    userID,
		"channel_id": channelID,
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	channel, ok := m.channels[channelID]
	if !ok {
		return
	}

	eventBytes, _ := models.MarshalJSONGeneric(event)
	for _, client := range channel {
		err := client.Conn.WriteMessage(1, eventBytes)
		if err != nil {
			return
		}
	}
}

// GetTokenForClient retrieves the token associated with a specific client connection.
// If needed, we identify the client by conn and channelID.
func (m *ConnectionManager) GetTokenForClient(conn *websocket.Conn, channelID int) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	channel, ok := m.channels[channelID]
	if !ok {
		return ""
	}
	client, ok := channel[conn]
	if !ok {
		return ""
	}
	return client.Token
}
