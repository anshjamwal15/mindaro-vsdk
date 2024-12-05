package repositories

import (
	"github.com/aimbot1526/mindaro-vsdk/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	Repository
}

// NewMessageRepository creates a new MessageRepository instance.
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{Repository{DB: db}}
}

// CreateMessage stores a new message in the database.
func (r *MessageRepository) CreateMessage(message *models.Message) (*models.Message, error) {
	if err := r.DB.Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// GetMessagesByGroup retrieves messages for a specific group.
func (r *MessageRepository) GetMessagesByGroup(groupID uint) ([]models.Message, error) {
	var messages []models.Message
	err := r.DB.Where("group_id = ?", groupID).Find(&messages).Error
	return messages, err
}

// MarkMessagesAsRead marks messages in a group as read.
func (r *MessageRepository) MarkMessagesAsRead(groupID uint, userID uint) error {
	return r.DB.Model(&models.Message{}).
		Where("group_id = ? AND receiver_id = ? AND is_read = ?", groupID, userID, false).
		Update("is_read", true).
		Error
}
