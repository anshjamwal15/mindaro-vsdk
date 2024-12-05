package repositories

import (
	"github.com/aimbot1526/mindaro-vsdk/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Repository{DB: db}}
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, userID).Error
	return &user, err
}

// GetUserByUsername retrieves a user by their username.
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
