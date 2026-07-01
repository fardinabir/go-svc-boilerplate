package user

import (
	"net/http"

	"github.com/fardinabir/go-svc-boilerplate/pkg/response"
)

// User domain error codes. Convention: U{Abbr}{HTTPStatus}
var (
	ErrNotFound      = &response.ErrorCode{Code: "UNF404", Status: http.StatusNotFound, Message: "User not found"}
	ErrAlreadyExists = &response.ErrorCode{Code: "UAE409", Status: http.StatusConflict, Message: "User already exists"}
)
