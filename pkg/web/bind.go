// Package web provides reusable, business-agnostic HTTP plumbing shared by all
// domains: request binding/validation, the response envelope, and the validator.
package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Base provides common request helpers for domain handlers to embed.
type Base struct{}

// MustBind binds the request body into req and validates it.
func (b Base) MustBind(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
