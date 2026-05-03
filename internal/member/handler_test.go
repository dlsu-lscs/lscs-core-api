package member

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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

// helper to create member info rows with 19 columns matching GetMemberInfo query (including image_url)
func createMemberInfoRow(email, fullName string) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "email", "full_name", "nickname", "image_url",
		"committee_id", "committee_name",
		"division_id", "division_name",
		"position_id", "position_name",
		"house_name",
		"contact_number", "college", "program",
		"interests", "discord", "fb_link", "telegram",
	}).AddRow(
		1, email, fullName, nil, nil,
		"RND", "Research and Development",
		"INT", "Internal",
		"MEM", "Member",
		"Gell-Mann",
		nil, "CCS", "CS-ST",
		nil, nil, nil, nil,
	)
}

// helper to create member info by ID rows with 19 columns matching GetMemberInfoById query (including image_url)
func createMemberInfoByIDRow(id int, email, fullName string) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "email", "full_name", "nickname", "image_url",
		"committee_id", "committee_name",
		"division_id", "division_name",
		"position_id", "position_name",
		"house_name",
		"contact_number", "college", "program",
		"interests", "discord", "fb_link", "telegram",
	}).AddRow(
		id, email, fullName, nil, nil,
		"RND", "Research and Development",
		"INT", "Internal",
		"MEM", "Member",
		"Gell-Mann",
		nil, "CCS", "CS-ST",
		nil, nil, nil, nil,
	)
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

		rows := createMemberInfoRow(reqBody.Email, "Test User")
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
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})

	t.Run("validation error - invalid email", func(t *testing.T) {
		e := echo.New()
		reqBody := EmailRequest{Email: "not-an-email"}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/member", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.GetMemberInfo(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestGetMemberInfoByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		reqBody := IDRequest{ID: 123}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/member-id", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := createMemberInfoByIDRow(123, "test@dlsu.edu.ph", "Test User")
		mock.ExpectQuery("SELECT (.+) FROM members m").WithArgs(int32(reqBody.ID)).WillReturnRows(rows)

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.GetMemberInfoByID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		e := echo.New()
		reqBody := IDRequest{ID: 123}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/member-id", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM members m").WithArgs(int32(reqBody.ID)).WillReturnError(sql.ErrNoRows)

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.GetMemberInfoByID(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})

	t.Run("validation error - invalid id", func(t *testing.T) {
		e := echo.New()
		reqBody := IDRequest{ID: 0} // invalid: must be > 0
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/member-id", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.GetMemberInfoByID(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestGetAllMembersHandler(t *testing.T) {
	t.Run("no filters", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/members", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id", "full_name", "nickname", "email", "telegram",
			"position_id", "committee_id", "college", "program",
			"discord", "interests", "contact_number", "fb_link", "image_url", "house_name",
		}).AddRow(
			1, "Test User 1", nil, "test1@dlsu.edu.ph", nil,
			"MEM", "RND", "CCS", "CS-ST",
			nil, nil, nil, nil, nil, "Gell-Mann",
		)
		mock.ExpectQuery("SELECT (.+) FROM members m").WithArgs("", "", "", "").WillReturnRows(rows)

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.GetAllMembersHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("with filters", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/members?position=MEM,EXE&committee=RND&committee=HR", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{
			"id", "full_name", "nickname", "email", "telegram",
			"position_id", "committee_id", "college", "program",
			"discord", "interests", "contact_number", "fb_link", "image_url", "house_name",
		}).AddRow(
			1, "Test User 1", nil, "test1@dlsu.edu.ph", nil,
			"MEM", "RND", "CCS", "CS-ST",
			nil, nil, nil, nil, nil, "Gell-Mann",
		)
		mock.ExpectQuery("SELECT (.+) FROM members m").WithArgs("MEM,EXE", "MEM,EXE", "RND,HR", "RND,HR").WillReturnRows(rows)

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.GetAllMembersHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
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
		reqBody := IDRequest{ID: 123}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/check-id", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM members WHERE id = ?").WithArgs(int32(reqBody.ID)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(reqBody.ID))

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.CheckIDIfMember(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("non-existing", func(t *testing.T) {
		e := echo.New()
		reqBody := IDRequest{ID: 123}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/check-id", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM members WHERE id = ?").WithArgs(int32(reqBody.ID)).WillReturnError(sql.ErrNoRows)

		dbService := &mockDBService{db: db}
		h := NewHandler(dbService)

		if assert.NoError(t, h.CheckIDIfMember(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
}
