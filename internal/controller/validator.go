package controller

import (
	"github.com/fardinabir/go-svc-boilerplate/internal/model"
	"github.com/go-playground/validator/v10"
)

// CustomValidator is a custom validator for the echo framework
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the input struct
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

// NewCustomValidator return a custom validator struct
func NewCustomValidator() *CustomValidator {
	v := validator.New()
	// Register the custom validation for user name using model function reference
	_ = v.RegisterValidation("validUserName", model.IsValidUserName)
	return &CustomValidator{validator: v}
}
