package repositories

import (
	"github.com/aimbot1526/mindaro-vsdk/models"
	"gorm.io/gorm"
)

type GroupRepository struct {
	Repository
}

// NewGroupRepository creates a new GroupRepository instance.
func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{Repository{DB: db}}
}

// CreateGroup creates a new group in the database.
func (r *GroupRepository) CreateGroup(group *models.Group) error {
	return r.DB.Create(group).Error
}

// AddMember adds a user to a group.
func (r *GroupRepository) AddMember(groupID uint, userID uint) error {
	groupMember := models.GroupMember{
		GroupID: groupID,
		UserID:  userID,
	}
	return r.DB.Create(&groupMember).Error
}

// GetGroupByID retrieves a group by its ID.
func (r *GroupRepository) GetGroupByID(groupID uint) (*models.Group, error) {
	var group models.Group
	err := r.DB.Preload("GroupMembers").First(&group, groupID).Error
	return &group, err
}

func (r *GroupRepository) GetGroupMembers(groupID uint) ([]models.User, error) {
	var members []models.User
	err := r.DB.Joins("JOIN group_members ON group_members.user_id = users.id").
		Where("group_members.group_id = ?", groupID).Find(&members).Error
	return members, err
}
