package handlers

import (
	"encoding/json"
	"github.com/genryusaishigikuni/messenger/gateway-service/internal/authclient"
	"github.com/genryusaishigikuni/messenger/gateway-service/internal/messageclient"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type IncomingMessage struct {
	ChannelID int    `json:"channel_id"`
	Content   string `json:"content"`
}

var upgraded = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(manager *ConnectionManager, authURL, messageURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Incoming WS connection request")
		token := r.URL.Query().Get("token")
		if token == "" {
			log.Println("No token provided")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		userID, err := authclient.ValidateToken(token, authURL)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("Token valid for userID: %d", userID)

		conn, err := upgraded.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		channelID := 1
		manager.RegisterClient(conn, userID, channelID, token)
		go handleClientMessages(conn, manager, messageURL, userID, channelID)
		log.Println("WebSocket connection established")
	}
}

func handleClientMessages(conn *websocket.Conn, manager *ConnectionManager, messageURL string, userID, channelID int) {
	defer manager.UnregisterClient(conn, channelID)

	// Retrieve token from manager
	token := manager.GetTokenForClient(conn, channelID)

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading ws message: %v", err)
			return
		}

		var incMsg IncomingMessage
		if err := json.Unmarshal(msgBytes, &incMsg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		storedMsg, err := messageclient.CreateMessage(messageURL, token, userID, incMsg.ChannelID, incMsg.Content)
		if err != nil {
			log.Printf("Failed to store message: %v", err)
			continue
		}

		manager.BroadcastToChannel(incMsg.ChannelID, storedMsg)
	}
}
