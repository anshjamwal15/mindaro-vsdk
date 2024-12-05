package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aimbot1526/mindaro-vsdk/models"
	"github.com/aimbot1526/mindaro-vsdk/repositories"
	"github.com/gorilla/mux"
)

type MessageHandler struct {
	MessageRepo *repositories.MessageRepository
}

func NewMessageHandler(messageRepo *repositories.MessageRepository) *MessageHandler {
	return &MessageHandler{MessageRepo: messageRepo}
}

// SendMessageToGroup allows a user to send a message to a group.
func (h *MessageHandler) SendMessageToGroup(w http.ResponseWriter, r *http.Request) {
	groupIDStr := mux.Vars(r)["group_id"]
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	var messageRequest struct {
		SenderID    uint   `json:"sender_id"`
		Content     string `json:"content"`
		MessageType string `json:"message_type"`
	}
	err = json.NewDecoder(r.Body).Decode(&messageRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message := &models.Message{
		SenderID:    messageRequest.SenderID,
		Content:     messageRequest.Content,
		MessageType: messageRequest.MessageType,
		GroupID:     uint(groupID),
	}

	createdMessage, err := h.MessageRepo.CreateMessage(message)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Message sent successfully",
		"message_id": createdMessage.ID,
	})
}
