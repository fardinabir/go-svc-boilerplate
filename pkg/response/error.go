package response

import (
	"crypto/rand"
	"fmt"
)

// ErrorCode describes an error class: a machine-readable code, HTTP status, and
// client-facing message. Declare instances as package-level vars:
//   - internal/errors/codes.go  — common cross-domain errors
//   - internal/<domain>/errors.go — domain-specific errors
//
// Code naming convention: {DomainAbbr}{ShortDesc}{HTTPStatus}
//
//	Common:  ISE500, BR400, NF404, UA401, F403, C409, UP422
//	User:    UNF404, UAE409
//	Cases:   CNF404, CANF404
//	Billing: BNF404, BDUP409
type ErrorCode struct {
	Code    string // machine-readable, namespaced
	Status  int    // HTTP status code
	Message string // client-facing message; never include internal detail here
}

func (e *ErrorCode) Error() string {
	return fmt.Sprintf("code=%s status=%d: %s", e.Code, e.Status, e.Message)
}

// NewID returns a short random hex string used as a correlation ID on every
// error response, so clients can reference it in support requests.
func NewID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
