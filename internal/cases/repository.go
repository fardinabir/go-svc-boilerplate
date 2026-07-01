package cases

import "gorm.io/gorm"

// Repository provides database operations for cases.
type Repository interface {
	Create(c *Case) error
	FindAll() ([]Case, error)
	FindByID(id int) (*Case, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new cases repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new case.
func (r *repository) Create(c *Case) error {
	return r.db.Create(c).Error
}

// FindAll retrieves all cases ordered by created_at desc.
func (r *repository) FindAll() ([]Case, error) {
	var cases []Case
	err := r.db.Order("created_at desc").Find(&cases).Error
	if err != nil {
		return nil, err
	}
	return cases, nil
}

// FindByID retrieves a case by ID.
func (r *repository) FindByID(id int) (*Case, error) {
	var c Case
	err := r.db.First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
