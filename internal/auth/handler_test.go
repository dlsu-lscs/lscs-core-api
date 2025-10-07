package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// mockAuthService is a mock implementation of the auth.Service interface.	ype mockAuthService struct{}

func (m *mockAuthService) GenerateJWT(email string) (string, error) {
	return "test_jwt_token", nil
}

// mockDBService is a mock implementation of the database.Service interface.
type mockDBService struct {
	db *sql.DB
}

func (m *mockDBService) Health() map[string]string {
	return nil
}

func (m *mockDBService) Close() error {
	return nil
}

func (m *mockDBService) GetConnection() *sql.DB {
	return m.db
}

func TestRequestKeyHandler(t *testing.T) {
	t.Run("success - RND member", func(t *testing.T) {
		e := echo.New()
		reqBody := EmailRequest{Email: "test@dlsu.edu.ph"}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/request-key", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM members m").
			WithArgs(reqBody.Email).
			WillReturnRows(sqlmock.NewRows([]string{"email", "full_name", "committee_name", "division_name", "position_name", "committee_id", "division_id", "house_name"}).
				AddRow(reqBody.Email, "Test User", "Research and Development", "", "", "", "", ""))

		mock.ExpectExec("INSERT INTO api_keys").
			WillReturnResult(sqlmock.NewResult(1, 1))

		dbService := &mockDBService{db: db}
		authService := &mockAuthService{}
		h := NewHandler(authService, dbService)

		if assert.NoError(t, h.RequestKeyHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var resp map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.Equal(t, reqBody.Email, resp["email"])
			assert.Equal(t, "test_jwt_token", resp["api_key"])
		}
	})

	t.Run("fail - non-RND member", func(t *testing.T) {
		e := echo.New()
		reqBody := EmailRequest{Email: "test@dlsu.edu.ph"}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/request-key", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM members m").
			WithArgs(reqBody.Email).
			WillReturnRows(sqlmock.NewRows([]string{"email", "full_name", "committee_name", "division_name", "position_name", "committee_id", "division_id", "house_name"}).
				AddRow(reqBody.Email, "Test User", "Not RND", "", "", "", "", ""))

		dbService := &mockDBService{db: db}
		authService := &mockAuthService{}
		h := NewHandler(authService, dbService)

		if assert.NoError(t, h.RequestKeyHandler(c)) {
			assert.Equal(t, http.StatusForbidden, rec.Code)
		}
	})

	t.Run("fail - not an LSCS member", func(t *testing.T) {
		e := echo.New()
		reqBody := EmailRequest{Email: "test@dlsu.edu.ph"}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/request-key", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM members m").
			WithArgs(reqBody.Email).
			WillReturnError(sql.ErrNoRows)

		dbService := &mockDBService{db: db}
		authService := &mockAuthService{}
		h := NewHandler(authService, dbService)

		if assert.NoError(t, h.RequestKeyHandler(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
}
