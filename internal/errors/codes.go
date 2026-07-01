package errors

import (
	"net/http"

	"github.com/fardinabir/go-svc-boilerplate/pkg/response"
)

// Common cross-domain error codes. Domain-specific errors live in
// internal/<domain>/errors.go and follow the same naming convention.
var (
	ErrInternalServerError = &response.ErrorCode{Code: "ISE500", Status: http.StatusInternalServerError, Message: "Internal server error"}
	ErrBadRequest          = &response.ErrorCode{Code: "BR400", Status: http.StatusBadRequest, Message: "Bad request"}
	ErrNotFound            = &response.ErrorCode{Code: "NF404", Status: http.StatusNotFound, Message: "Resource not found"}
	ErrUnauthorized        = &response.ErrorCode{Code: "UA401", Status: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden           = &response.ErrorCode{Code: "F403", Status: http.StatusForbidden, Message: "Forbidden"}
	ErrConflict            = &response.ErrorCode{Code: "C409", Status: http.StatusConflict, Message: "Conflict"}
	ErrUnprocessable       = &response.ErrorCode{Code: "UP422", Status: http.StatusUnprocessableEntity, Message: "Unprocessable entity"}
)
