package member

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/labstack/echo/v4"
)

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

func TestGetMemberInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		reqBody := EmailRequest{Email: "test@dlsu.edu.ph"}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/member", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"email", "full_name", "committee_name", "division_name", "position_name", "committee_id", "division_id", "house_name"}).
			AddRow(reqBody.Email, "Test User", "RND", "INT", "Trainee", "RND", "INT", "Gell-Mann")
		mock.ExpectQuery("SELECT (.+) FROM members m").WithArgs(reqBody.Email).WillReturnRows(rows)

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.GetMemberInfo(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		e := echo.New()
		reqBody := EmailRequest{Email: "test@dlsu.edu.ph"}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/member", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM members m").WithArgs(reqBody.Email).WillReturnError(sql.ErrNoRows)

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.GetMemberInfo(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestGetMemberInfoByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		reqBody := IdRequest{Id: 123}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/member-id", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "email", "full_name", "committee_name", "division_name", "position_name", "committee_id", "division_id", "house_name"}).
			AddRow(123, "test@dlsu.edu.ph", "Test User", "RND", "INT", "Trainee", "RND", "INT", "Gell-Mann")
		mock.ExpectQuery("SELECT (.+) FROM members m").WithArgs(int32(reqBody.Id)).WillReturnRows(rows)

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.GetMemberInfoByID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		e := echo.New()
		reqBody := IdRequest{Id: 123}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/member-id", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM members m").WithArgs(int32(reqBody.Id)).WillReturnError(sql.ErrNoRows)

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.GetMemberInfoByID(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestGetAllMembersHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/members", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "full_name", "nickname", "email", "telegram", "position_id", "committee_id", "college", "program", "discord", "interests", "contact_number", "fb_link", "house_name"}).
		AddRow(1, "Test User 1", "Test1", "test1@dlsu.edu.ph", "", "MEM", "RND", "CCS", "CS-ST", "", "", "", "", "Gell-Mann")
	mock.ExpectQuery("SELECT (.+) FROM members m").WillReturnRows(rows)

	dbService := &mockDBService{db: db}
	h := NewHandler(dbService)

	if assert.NoError(t, h.GetAllMembersHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCheckEmailHandler(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		e := echo.New()
		reqBody := EmailRequest{Email: "test@dlsu.edu.ph"}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/check-email", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT email FROM members WHERE email = ?").WithArgs(reqBody.Email).WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow(reqBody.Email))

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.CheckEmailHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("non-existing", func(t *testing.T) {
		e := echo.New()
		reqBody := EmailRequest{Email: "test@dlsu.edu.ph"}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/check-email", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT email FROM members WHERE email = ?").WithArgs(reqBody.Email).WillReturnError(sql.ErrNoRows)

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.CheckEmailHandler(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
}

func TestCheckIDIfMember(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		e := echo.New()
		reqBody := IdRequest{Id: 123}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/check-id", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM members WHERE id = ?").WithArgs(int32(reqBody.Id)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(reqBody.Id))

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.CheckIDIfMember(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("non-existing", func(t *testing.T) {
		e := echo.New()
		reqBody := IdRequest{Id: 123}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/check-id", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM members WHERE id = ?").WithArgs(int32(reqBody.Id)).WillReturnError(sql.ErrNoRows)

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.CheckIDIfMember(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
}