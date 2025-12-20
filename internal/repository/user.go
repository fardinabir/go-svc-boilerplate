package repository

import (
	"github.com/fardinabir/go-svc-boilerplate/internal/model"
	"gorm.io/gorm"
)

// UserRepository provides database operations for users
type UserRepository interface {
	Create(user *model.User) error
	FindAll() ([]model.User, error)
	FindByID(id int) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindAll retrieves all users ordered by created_at desc
func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Order("created_at desc").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindByID retrieves a user by ID
func (r *userRepository) FindByID(id int) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
