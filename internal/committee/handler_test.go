package committee

import (
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

func TestGetAllCommitteesHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/committees", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"committee_id", "committee_name", "committee_head", "division_id"}).
		AddRow("RND", "Research and Development", 1, "INT")
	mock.ExpectQuery("SELECT (.+) FROM committees").WillReturnRows(rows)

	dbService := &mockDBService{db: db}
	h := NewHandler(dbService)

	if assert.NoError(t, h.GetAllCommitteesHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string][]repository.Committee
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Len(t, resp["committees"], 1)
		assert.Equal(t, "RND", resp["committees"][0].CommitteeID)
	}
}
