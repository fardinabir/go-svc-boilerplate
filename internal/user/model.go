// Package user is the user bounded context: it owns the User model and its
// repository, service, and HTTP handler. It is the reference template every other
// domain package copies.
package user

import (
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// User represents a simple user entity for boilerplate CRUD.
type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// RegisterValidations registers this domain's custom validation tags. The
// composition root passes it to web.NewCustomValidator, so the web package never
// needs to import this domain.
func RegisterValidations(v *validator.Validate) {
	_ = v.RegisterValidation("validUserName", IsValidUserName)
}

// IsValidUserName is a validator function compatible with go-playground/validator.
func IsValidUserName(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	trimmed := strings.TrimSpace(s)
	if len(trimmed) < 2 || len(trimmed) > 50 {
		return false
	}
	for _, r := range trimmed {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) && r != '-' && r != '\'' {
			return false
		}
	}
	return true
}
