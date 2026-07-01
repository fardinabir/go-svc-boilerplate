package web

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator is a custom validator for the echo framework.
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the input struct.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// NewCustomValidator builds a validator and applies each registrar so domains can
// register their own custom tags without this package importing the domains.
func NewCustomValidator(registrars ...func(*validator.Validate)) *CustomValidator {
	v := validator.New()
	for _, register := range registrars {
		register(v)
	}
	return &CustomValidator{validator: v}
}
