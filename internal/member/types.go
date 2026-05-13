package member

import (
	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/storage"
)

type Handler struct {
	dbService database.Service
	s3Service *storage.S3Service
}

func NewHandler(dbService database.Service, s3Service *storage.S3Service) *Handler {
	return &Handler{
		dbService: dbService,
		s3Service: s3Service,
	}

}
