package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	db2 "github.com/fardinabir/go-svc-boilerplate/internal/db"
	"github.com/fardinabir/go-svc-boilerplate/internal/model"
	"github.com/fardinabir/go-svc-boilerplate/internal/repository"
	"github.com/fardinabir/go-svc-boilerplate/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserHandler_CreateUser(t *testing.T) {
	e := echo.New()
	e.Validator = NewCustomValidator()
	dbInstance, err := db2.NewTestDB()
	require.NoError(t, err)
	err = db2.Migrate(dbInstance)
	require.NoError(t, err)
	repository := repository.NewUserRepository(dbInstance)
	service := service.NewUserService(repository)
	handler := NewUserHandler(service)

	tests := []struct {
		name       string
		createBody string
		statusCode int
	}{
		{name: "valid_create_user", createBody: `{"name":"Alice","email":"alice@example.com"}`, statusCode: http.StatusCreated},
		{name: "invalid_missing_email", createBody: `{"name":"Alice"}`, statusCode: http.StatusBadRequest},
		{name: "invalid_email_format", createBody: `{"name":"Alice","email":"not-an-email"}`, statusCode: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearDB(dbInstance, model.User{})
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(tt.createBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/users")
			err := handler.CreateUser(c)
			require.NoError(t, err)
			assert.Equal(t, tt.statusCode, rec.Code)
		})
	}
}

func TestUserHandler_ListAndGetUser(t *testing.T) {
	e := echo.New()
	e.Validator = NewCustomValidator()
	dbInstance, err := db2.NewTestDB()
	require.NoError(t, err)
	err = db2.Migrate(dbInstance)
	require.NoError(t, err)
	repository := repository.NewUserRepository(dbInstance)
	service := service.NewUserService(repository)
	handler := NewUserHandler(service)

	// Initially empty list
	{
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users")
		require.NoError(t, handler.ListUsers(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// Create a user
	{
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`{"name":"Bob","email":"bob@example.com"}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users")
		require.NoError(t, handler.CreateUser(c))
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// List again should have one user
	{
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users")
		require.NoError(t, handler.ListUsers(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// Get by invalid id
	{
		req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("abc")
		require.NoError(t, handler.GetUserByID(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

// Helpers
func clearDB(db *gorm.DB, models ...interface{}) {
	for _, model := range models {
		db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model)
	}
}
