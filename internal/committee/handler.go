package committee

import (
	"fmt"
	"net/http"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	dbService database.Service
}

func NewHandler(dbService database.Service) *Handler {
	return &Handler{
		dbService: dbService,
	}
}

func (h *Handler) GetAllCommitteesHandler(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	committees, err := q.GetAllCommittees(ctx)
	if err != nil {
		response := map[string]string{
			"error": fmt.Sprintf("Internal server error: %v", err),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"committees": committees,
	})
}
