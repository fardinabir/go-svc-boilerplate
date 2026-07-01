package cases

// Service provides case operations.
type Service interface {
	CreateCase(c *Case) error
	ListCases() ([]Case, error)
	GetCaseByID(id int) (*Case, error)
	// AssigneeEmail reads across the user domain through the UserReader port.
	AssigneeEmail(caseID int) (string, error)
}

type service struct {
	repo  Repository
	users UserReader
}

// NewService creates a new cases service. users is the cross-domain port, supplied
// by the composition root (user.Service satisfies it).
func NewService(repo Repository, users UserReader) Service {
	return &service{repo: repo, users: users}
}

// CreateCase creates a new case.
func (s *service) CreateCase(c *Case) error {
	return s.repo.Create(c)
}

// ListCases retrieves all cases.
func (s *service) ListCases() ([]Case, error) {
	return s.repo.FindAll()
}

// GetCaseByID retrieves a single case by id.
func (s *service) GetCaseByID(id int) (*Case, error) {
	return s.repo.FindByID(id)
}

// AssigneeEmail returns the email of the user assigned to the case, fetched via the
// user domain's port — without cases knowing anything about the user model.
func (s *service) AssigneeEmail(caseID int) (string, error) {
	c, err := s.repo.FindByID(caseID)
	if err != nil {
		return "", err
	}
	return s.users.EmailByID(c.AssigneeID)
}
