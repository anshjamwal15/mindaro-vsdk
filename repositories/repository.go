package repositories

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new repository with the given database connection.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}
