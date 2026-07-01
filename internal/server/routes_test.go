package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fardinabir/go-svc-boilerplate/internal/cases"
	"github.com/fardinabir/go-svc-boilerplate/internal/db"
	"github.com/fardinabir/go-svc-boilerplate/internal/health"
	"github.com/fardinabir/go-svc-boilerplate/internal/user"
	"github.com/fardinabir/go-svc-boilerplate/pkg/web"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRegisterRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	dbInstance, err := db.NewTestDB()
	require.NoError(t, err)
	require.NoError(t, db.Migrate(dbInstance))
	require.NoError(t, setupTestRoutes(e, dbInstance))

	// Test cases
	tests := []struct {
		name         string
		method       string
		target       string
		expectedCode int
	}{
		{"Health_Check", http.MethodGet, "/api/v1/health", http.StatusOK},
		{"Create_User_without_body", http.MethodPost, "/api/v1/users", http.StatusBadRequest},
		{"List_users_empty", http.MethodGet, "/api/v1/users", http.StatusOK},
		{"Create_Case_without_body", http.MethodPost, "/api/v1/cases", http.StatusBadRequest},
		{"List_cases_empty", http.MethodGet, "/api/v1/cases", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.target, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

// setupTestRoutes mirrors the production wiring in setupRoutes: build handlers via
// Wire, then register each domain's routes.
func setupTestRoutes(e *echo.Echo, dbInstance *gorm.DB) error {
	e.Validator = web.NewCustomValidator(user.RegisterValidations)

	handlers, err := InitializeHandlers(dbInstance)
	if err != nil {
		return err
	}

	api := e.Group("/api/v1")
	health.RegisterRoutes(api, health.New())
	user.RegisterRoutes(api.Group("/users"), handlers.User)
	cases.RegisterRoutes(api.Group("/cases"), handlers.Cases)
	return nil
}
