package handlers

import (
	"strconv"
	"sync"

	"github.com/genryusaishigikuni/messenger/gateway-service/pkg/models"
	"github.com/genryusaishigikuni/messenger/gateway-service/pkg/utils"
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
	utils.Info("Initializing Connection Manager")
	return &ConnectionManager{
		channels: make(map[int]map[*websocket.Conn]*clientInfo),
	}
}

func (m *ConnectionManager) RegisterClient(conn *websocket.Conn, userID, channelID int, token string) {
	utils.Info("Registering new client")
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.channels[channelID]; !ok {
		utils.Info("Creating a new channel for channel ID: " + strconv.Itoa(channelID))
		m.channels[channelID] = make(map[*websocket.Conn]*clientInfo)
	}
	m.channels[channelID][conn] = &clientInfo{
		Conn:      conn,
		UserID:    userID,
		ChannelID: channelID,
		Token:     token,
	}
	utils.Info("Client registered: UserID=" + strconv.Itoa(userID) + ", ChannelID=" + strconv.Itoa(channelID))
}

func (m *ConnectionManager) UnregisterClient(conn *websocket.Conn, channelID int) {
	utils.Info("Unregistering client")
	m.mu.Lock()
	defer m.mu.Unlock()

	if channel, ok := m.channels[channelID]; ok {
		delete(channel, conn)
		if len(channel) == 0 {
			utils.Info("Channel is empty; deleting channel ID: " + strconv.Itoa(channelID))
			delete(m.channels, channelID)
		}
	}

	err := conn.Close()
	if err != nil {
		utils.Error("Failed to close WebSocket connection: " + err.Error())
	} else {
		utils.Info("WebSocket connection closed")
	}
}

func (m *ConnectionManager) BroadcastToChannel(channelID int, msg models.Message) {
	utils.Info("Broadcasting message to channel: " + strconv.Itoa(channelID))
	m.mu.RLock()
	defer m.mu.RUnlock()

	channel, ok := m.channels[channelID]
	if !ok {
		utils.Error("Channel not found for ID: " + strconv.Itoa(channelID))
		return
	}

	msgBytes, err := msg.MarshalJSON()
	if err != nil {
		utils.Error("Failed to marshal message: " + err.Error())
		return
	}

	for _, client := range channel {
		err := client.Conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			utils.Error("Failed to send message to client: " + err.Error())
		}
	}
	utils.Info("Message broadCasted successfully to channel: " + strconv.Itoa(channelID))
}

func (m *ConnectionManager) BroadcastPresenceEvent(eventType string, userID, channelID int) {
	utils.Info("Broadcasting presence event: " + eventType)
	event := map[string]interface{}{
		"event":      eventType,
		"user_id":    userID,
		"channel_id": channelID,
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	channel, ok := m.channels[channelID]
	if !ok {
		utils.Error("Channel not found for ID: " + strconv.Itoa(channelID))
		return
	}

	eventBytes, err := models.MarshalJSONGeneric(event)
	if err != nil {
		utils.Error("Failed to marshal presence event: " + err.Error())
		return
	}

	for _, client := range channel {
		err := client.Conn.WriteMessage(websocket.TextMessage, eventBytes)
		if err != nil {
			utils.Error("Failed to send presence event to client: " + err.Error())
		}
	}
	utils.Info("Presence event broadCasted successfully to channel: " + strconv.Itoa(channelID))
}

func (m *ConnectionManager) GetTokenForClient(conn *websocket.Conn, channelID int) string {
	utils.Info("Fetching token for client in channel ID: " + strconv.Itoa(channelID))
	m.mu.RLock()
	defer m.mu.RUnlock()

	channel, ok := m.channels[channelID]
	if !ok {
		utils.Error("Channel not found for ID: " + strconv.Itoa(channelID))
		return ""
	}
	client, ok := channel[conn]
	if !ok {
		utils.Error("Client connection not found in channel ID: " + strconv.Itoa(channelID))
		return ""
	}
	utils.Info("Token fetched successfully for client")
	return client.Token
}
