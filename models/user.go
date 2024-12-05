package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `json:"username" gorm:"unique"`
	Password  string `json:"password" gorm:"not null"`
	Email     string `json:"email" gorm:"unique"`
	IsOnline  bool   `json:"is_online"`
	Blocked   bool   `json:"blocked"`
	SessionID string `json:"session_id"`
}
