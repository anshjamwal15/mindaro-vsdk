package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	SessionID string `json:"session_id" gorm:"unique"`
}
