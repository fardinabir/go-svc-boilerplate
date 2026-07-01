package cases_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fardinabir/go-svc-boilerplate/internal/cases"
	"github.com/fardinabir/go-svc-boilerplate/internal/db"
	"github.com/fardinabir/go-svc-boilerplate/pkg/web"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// stubUserReader satisfies cases.UserReader without touching the real user domain,
// showing the cross-domain port is mockable in isolation.
type stubUserReader struct{ email string }

func (s stubUserReader) EmailByID(_ int) (string, error) { return s.email, nil }

func newCasesHandler(t *testing.T) (cases.Handler, *gorm.DB) {
	t.Helper()
	dbInstance, err := db.NewTestDB()
	require.NoError(t, err)
	require.NoError(t, db.Migrate(dbInstance))
	svc := cases.NewService(cases.NewRepository(dbInstance), stubUserReader{email: "assignee@example.com"})
	return cases.NewHandler(svc), dbInstance
}

func TestCaseHandler_CreateListGet(t *testing.T) {
	e := echo.New()
	e.Validator = web.NewCustomValidator()
	handler, dbInstance := newCasesHandler(t)
	clearTables(dbInstance, "cases")

	// Create
	{
		body := `{"file_number":"FC-1001","status":"open","servicer_id":1,"assignee_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/cases")
		require.NoError(t, handler.CreateCase(c))
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// List
	{
		req := httptest.NewRequest(http.MethodGet, "/cases", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/cases")
		require.NoError(t, handler.ListCases(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// Get by invalid id
	{
		req := httptest.NewRequest(http.MethodGet, "/cases/abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/cases/:id")
		c.SetParamNames("id")
		c.SetParamValues("abc")
		require.NoError(t, handler.GetCaseByID(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

// clearTables deletes rows from tables in the given order.
// Caller must list FK children before parents (e.g. "cases" before "users").
func clearTables(dbInstance *gorm.DB, tables ...string) {
	for _, t := range tables {
		dbInstance.Exec(`DELETE FROM "` + t + `"`)
	}
}
