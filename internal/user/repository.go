package user

import (
	"gorm.io/gorm"
)

// Repository provides database operations for users.
type Repository interface {
	Create(user *User) error
	FindAll() ([]User, error)
	FindByID(id int) (*User, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new user.
func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

// FindAll retrieves all users ordered by created_at desc.
func (r *repository) FindAll() ([]User, error) {
	var users []User
	err := r.db.Order("created_at desc").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindByID retrieves a user by ID.
func (r *repository) FindByID(id int) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
