package models

import (
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// Group represents a chat group
type Group struct {
	gorm.Model
	GroupID      string `gorm:"uniqueIndex;size:10" json:"group_id"` // Unique alphanumeric ID for the group
	Name         string `json:"name"`                                // Name of the group
	IsPrivate    bool   `json:"is_private"`                          // Whether the group is private (2 users) or public
	CreatorID    uint   `json:"creator_id"`                          // User ID of the creator
	GroupMembers []User `gorm:"many2many:group_members;" json:"group_members"`
}

// GroupMember represents the association between a user and a group
type GroupMember struct {
	gorm.Model
	GroupID uint `gorm:"index"`
	UserID  uint `gorm:"index"`
}

// generateRandomID creates a random alphanumeric string of given length
func generateRandomID(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a local random generator
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, length)
	for i := range id {
		id[i] = charset[r.Intn(len(charset))]
	}
	return string(id)
}

// BeforeCreate hook to generate a unique GroupID
func (g *Group) BeforeCreate(tx *gorm.DB) (err error) {
	g.GroupID = generateRandomID(10) // Generate a 10-character random ID
	return
}
