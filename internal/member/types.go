package member

import (
	"github.com/dlsu-lscs/lscs-core-api/internal/database"
)

type Handler struct {
	dbService database.Service
}

func NewHandler(dbService database.Service) *Handler {
	return &Handler{
		dbService: dbService,
	}

}
