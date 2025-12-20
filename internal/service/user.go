package service

import (
	"github.com/fardinabir/go-svc-boilerplate/internal/model"
	"github.com/fardinabir/go-svc-boilerplate/internal/repository"
)

// UserService provides user operations
type UserService interface {
	CreateUser(user *model.User) error
	ListUsers() ([]model.User, error)
	GetUserByID(id int) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// CreateUser creates a new user
func (s *userService) CreateUser(user *model.User) error {
	return s.repo.Create(user)
}

// ListUsers retrieves all users
func (s *userService) ListUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

// GetUserByID retrieves a single user by id
func (s *userService) GetUserByID(id int) (*model.User, error) {
	return s.repo.FindByID(id)
}
