package user

// Service provides user operations.
type Service interface {
	CreateUser(user *User) error
	ListUsers() ([]User, error)
	GetUserByID(id int) (*User, error)
	// EmailByID exposes a user's email to other domains; it satisfies
	// cases.UserReader (bound in the composition root).
	EmailByID(id int) (string, error)
}

type service struct {
	repo Repository
}

// NewService creates a new user service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateUser creates a new user.
func (s *service) CreateUser(user *User) error {
	return s.repo.Create(user)
}

// ListUsers retrieves all users.
func (s *service) ListUsers() ([]User, error) {
	return s.repo.FindAll()
}

// GetUserByID retrieves a single user by id.
func (s *service) GetUserByID(id int) (*User, error) {
	return s.repo.FindByID(id)
}

// EmailByID returns the email of the user with the given id.
func (s *service) EmailByID(id int) (string, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
