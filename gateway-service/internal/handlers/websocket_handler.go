package handlers

import (
	"encoding/json"
	"github.com/genryusaishigikuni/messenger/gateway-service/internal/authclient"
	"github.com/genryusaishigikuni/messenger/gateway-service/internal/messageclient"
	"github.com/genryusaishigikuni/messenger/gateway-service/pkg/utils"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
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
		utils.Info("Incoming WebSocket connection request")

		// Extract token from query
		token := r.URL.Query().Get("token")
		if token == "" {
			utils.Error("No token provided")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Validate token with auth service
		userID, err := authclient.ValidateToken(token, authURL)
		if err != nil {
			utils.Error("Token validation failed: " + err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		utils.Info("Token validated successfully for userID: " + strconv.Itoa(userID))

		// Upgrade the connection to WebSocket
		conn, err := upgraded.Upgrade(w, r, nil)
		if err != nil {
			utils.Error("WebSocket upgrade error: " + err.Error())
			return
		}

		utils.Info("WebSocket connection established")

		// Assign default channel ID for now
		channelID := 1
		manager.RegisterClient(conn, userID, channelID, token)
		utils.Info("Client registered: UserID=" + strconv.Itoa(userID) + ", ChannelID=" + strconv.Itoa(channelID))

		go handleClientMessages(conn, manager, messageURL, userID, channelID)
	}
}

func handleClientMessages(conn *websocket.Conn, manager *ConnectionManager, messageURL string, userID, channelID int) {
	defer func() {
		utils.Info("Unregistering client: UserID=" + strconv.Itoa(userID) + ", ChannelID=" + strconv.Itoa(channelID))
		manager.UnregisterClient(conn, channelID)
	}()

	// Retrieve token from manager
	token := manager.GetTokenForClient(conn, channelID)

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			utils.Error("Error reading WebSocket message: " + err.Error())
			return
		}

		// Parse the incoming message
		var incMsg IncomingMessage
		if err := json.Unmarshal(msgBytes, &incMsg); err != nil {
			utils.Error("Invalid message format: " + err.Error())
			continue
		}

		utils.Info("Received message: ChannelID=" + strconv.Itoa(incMsg.ChannelID) + ", Content=" + incMsg.Content)

		// Create and store message using the message service
		storedMsg, err := messageclient.CreateMessage(messageURL, token, userID, incMsg.ChannelID, incMsg.Content)
		if err != nil {
			utils.Error("Failed to store message: " + err.Error())
			continue
		}

		// Broadcast the stored message to the channel
		manager.BroadcastToChannel(incMsg.ChannelID, storedMsg)
		utils.Info("Message broadCasted: ChannelID=" + strconv.Itoa(incMsg.ChannelID))
	}
}
