//go:build tools

// Package tools pins build-time tool dependencies (the Wire CLI) so `go mod tidy`
// keeps them in go.mod. It is never compiled into the application binary.
package tools

import (
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/swaggo/swag/cmd/swag"
)
