package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID    uint   `json:"sender_id"`   // User ID of the sender
	ReceiverID  uint   `json:"receiver_id"` // Used for private messages (could be a group ID for group messages)
	Content     string `json:"content"`
	IsRead      bool   `json:"is_read"`            // Track if the message was read
	MessageType string `json:"message_type"`       // e.g., "text", "image", etc.
	GroupID     uint   `json:"group_id,omitempty"` // For group messages, store the group ID
}
