package cases

import (
	"net/http"

	"github.com/fardinabir/go-svc-boilerplate/pkg/response"
)

// Cases domain error codes. Convention: C{Abbr}{HTTPStatus}
var (
	ErrNotFound         = &response.ErrorCode{Code: "CNF404", Status: http.StatusNotFound, Message: "Case not found"}
	ErrAssigneeNotFound = &response.ErrorCode{Code: "CANF404", Status: http.StatusNotFound, Message: "Case assignee not found"}
)
