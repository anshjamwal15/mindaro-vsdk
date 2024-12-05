package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/aimbot1526/mindaro-vsdk/models"
	"github.com/aimbot1526/mindaro-vsdk/repositories"
	"github.com/aimbot1526/mindaro-vsdk/utils"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	GroupRepo   *repositories.GroupRepository
	MessageRepo *repositories.MessageRepository
	Clients     map[string]map[*websocket.Conn]bool // groupID -> connections
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewWebSocketHandler(groupRepo *repositories.GroupRepository, messageRepo *repositories.MessageRepository) *WebSocketHandler {
	return &WebSocketHandler{
		GroupRepo:   groupRepo,
		MessageRepo: messageRepo,
		Clients:     make(map[string]map[*websocket.Conn]bool),
	}
}

// WebSocket handler
func (h *WebSocketHandler) GroupWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	clientAddress := r.RemoteAddr
	log.Printf("Client connected: %s", clientAddress)

	// Extract group_id from URL query
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		conn.WriteJSON(map[string]string{"error": "group_id is required"})
		return
	}

	// Register connection
	if _, exists := h.Clients[groupID]; !exists {
		h.Clients[groupID] = make(map[*websocket.Conn]bool)
	}
	h.Clients[groupID][conn] = true

	defer func() {
		// Clean up connection when closed
		delete(h.Clients[groupID], conn)
		conn.Close()
	}()

	// Listen for incoming messages
	for {
		var msg struct {
			SenderID    string `json:"sender_id"`
			Content     string `json:"content"`
			MessageType string `json:"message_type"`
		}

		// Read the incoming message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)

			// Handle abnormal closure gracefully
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("Connection closed gracefully")
			} else {
				log.Println("Unexpected WebSocket error:", err)
			}
			break
		}

		// Log the message
		log.Printf("Received message: %+v", msg)

		// Process message (save to DB and broadcast)
		message := &models.Message{
			GroupID:     utils.ParseStringToUint(groupID),
			SenderID:    utils.ParseStringToUint(msg.SenderID),
			Content:     msg.Content,
			MessageType: msg.MessageType,
		}

		// Save the message to the database
		messageID, err := h.MessageRepo.CreateMessage(message)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			conn.WriteJSON(map[string]string{"error": "Failed to save message"})
			continue
		}

		// Broadcast the message to all connections in the group
		response := map[string]interface{}{
			"group_id":     groupID,
			"message_id":   messageID,
			"sender_id":    msg.SenderID,
			"content":      msg.Content,
			"message_type": msg.MessageType,
			"timestamp":    time.Now(), // Replace with actual timestamp
		}

		for client := range h.Clients[groupID] {
			err := client.WriteJSON(response)
			if err != nil {
				log.Printf("Error broadcasting message: %v", err)
				client.Close()
				delete(h.Clients[groupID], client)
			}
		}
	}
}

// Notify online users in the group when a member comes online or goes offline
func (h *WebSocketHandler) NotifyUserStatus(groupID, userID, status string) {
	notification := map[string]string{
		"user_id": userID,
		"status":  status, // e.g., online or offline
	}

	for client := range h.Clients[groupID] {
		err := client.WriteJSON(notification)
		if err != nil {
			log.Printf("Error sending user status notification: %v", err)
			client.Close()
			delete(h.Clients[groupID], client)
		}
	}
}
