package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	db2 "github.com/fardinabir/go-svc-boilerplate/internal/db"
	"github.com/fardinabir/go-svc-boilerplate/internal/repository"
	"github.com/fardinabir/go-svc-boilerplate/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRegister(t *testing.T) {
	// Setup
	e := echo.New()
	dbInstance, err := db2.NewTestDB()
	require.NoError(t, err)
	err = db2.Migrate(dbInstance)
	require.NoError(t, err)
	setupTestRoutes(e, dbInstance)

	// Test cases
	tests := []struct {
		name         string
		method       string
		target       string
		expectedCode int
	}{
		{"Health_Check", http.MethodGet, "/api/v1/health", http.StatusOK},
		{"Create_User_without_body", http.MethodPost, "/api/v1/users", http.StatusBadRequest}, // Assuming no body is sent, should return BadRequest
		{"List_users_empty", http.MethodGet, "/api/v1/users", http.StatusOK},                  // Should return empty array
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

// setupTestRoutes configures routes for testing with the same pattern as the server
func setupTestRoutes(e *echo.Echo, db *gorm.DB) {
	// Set up request validation
	e.Validator = NewCustomValidator()

	// Create API version group
	api := e.Group("/api/v1")

	// Register health check endpoint
	healthHandler := NewHealth()
	api.GET("/health", healthHandler.Health)

	// Initialize user handler with dependencies
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	// Register user routes
	InitRoutes(api, userHandler)
}
