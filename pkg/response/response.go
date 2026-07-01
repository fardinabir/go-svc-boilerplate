package response

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

// ResponseData is the generic success envelope.
type ResponseData struct {
	Data interface{} `json:"data,omitempty"`
}

// APIError is the error response body.
//
//	{
//	  "id":     "a3f2c1d4e5b6a7c8",   ← correlation ID; reference in support requests
//	  "code":   "CNF404",              ← machine-readable domain error code
//	  "title":  "Case not found",      ← client-facing message
//	  "detail": {...}                  ← optional structured detail (e.g. validation errors)
//	}
//
// Internal error causes (DB messages, stack traces) are never included.
// Log them server-side before calling Respond.
type APIError struct {
	ID     string          `json:"id"`
	Code   string          `json:"code"`
	Title  string          `json:"title"`
	Detail json.RawMessage `json:"detail,omitempty" swaggertype:"object"`
}

// Respond writes an APIError response from an ErrorCode.
// Pass an optional detail value (e.g. a validation error string or map) as the
// third argument — marshalled to JSON and included only if non-nil.
//
// Usage:
//
//	return response.Respond(c, ErrNotFound)
//	return response.Respond(c, apierr.ErrBadRequest, validationErrs)
func Respond(c echo.Context, ec *ErrorCode, detail ...interface{}) error {
	apiErr := &APIError{
		ID:    NewID(),
		Code:  ec.Code,
		Title: ec.Message,
	}
	if len(detail) > 0 && detail[0] != nil {
		if raw, err := json.Marshal(detail[0]); err == nil && string(raw) != "null" {
			apiErr.Detail = raw
		}
	}
	return c.JSON(ec.Status, apiErr)
}
