package user_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fardinabir/go-svc-boilerplate/internal/db"
	"github.com/fardinabir/go-svc-boilerplate/internal/user"
	"github.com/fardinabir/go-svc-boilerplate/pkg/web"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newUserHandler(t *testing.T) (user.Handler, *gorm.DB) {
	t.Helper()
	dbInstance, err := db.NewTestDB()
	require.NoError(t, err)
	require.NoError(t, db.Migrate(dbInstance))
	return user.NewHandler(user.NewService(user.NewRepository(dbInstance))), dbInstance
}

func TestUserHandler_CreateUser(t *testing.T) {
	e := echo.New()
	e.Validator = web.NewCustomValidator(user.RegisterValidations)
	handler, dbInstance := newUserHandler(t)

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
			clearDB(dbInstance, user.User{})
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
	e.Validator = web.NewCustomValidator(user.RegisterValidations)
	handler, dbInstance := newUserHandler(t)
	clearDB(dbInstance, user.User{})

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

// clearDB truncates the given models between test cases.
func clearDB(dbInstance *gorm.DB, models ...interface{}) {
	for _, m := range models {
		dbInstance.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(m)
	}
}
