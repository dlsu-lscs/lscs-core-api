package storage

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/config"
	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
)

// UploadHandler handles profile image upload requests
type UploadHandler struct {
	s3Service *S3Service
	dbService database.Service
	cfg       *config.Config
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(s3Service *S3Service, dbService database.Service, cfg *config.Config) *UploadHandler {
	return &UploadHandler{
		s3Service: s3Service,
		dbService: dbService,
		cfg:       cfg,
	}
}

// GenerateUploadURLRequest represents a request to generate an upload URL
type GenerateUploadURLRequest struct {
	ContentType string `json:"content_type" validate:"required,oneof=image/jpeg image/png image/webp"`
}

// GenerateUploadURLResponse represents the response for an upload URL request
type GenerateUploadURLResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
	ExpiresIn int    `json:"expires_in_seconds"`
}

// CompleteUploadRequest represents a request to complete an upload
type CompleteUploadRequest struct {
	ObjectKey string `json:"object_key" validate:"required"`
}

// CompleteUploadResponse represents the response for completing an upload
type CompleteUploadResponse struct {
	ImageURL    string `json:"image_url"`
	DownloadURL string `json:"download_url,omitempty"`
	Message     string `json:"message"`
}

// GenerateUploadURLHandler generates a pre-signed URL for uploading a profile image
// @Summary Generate upload URL for profile image
// @Description Generate a pre-signed URL for uploading a profile image
// @Tags upload
// @Accept json
// @Produce json
// @Param request body GenerateUploadURLRequest true "Upload request"
// @Success 200 {object} GenerateUploadURLResponse "Upload URL generated"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security SessionAuth
// @Router /upload/profile-image [post]
func (h *UploadHandler) GenerateUploadURLHandler(c echo.Context) error {
	if !h.s3Service.IsEnabled() {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Storage service not available",
		})
	}

	req := new(GenerateUploadURLRequest)
	if err := helpers.BindAndValidate(c, req); err != nil {
		return err
	}

	memberID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	log.Info().
		Int32("member_id", memberID).
		Str("content_type", req.ContentType).
		Msg("upload URL requested")

	uploadURL, objectKey, err := h.s3Service.GenerateUploadURL(c.Request().Context(), memberID, req.ContentType)
	if err != nil {
		log.Error().Err(err).Int32("member_id", memberID).Msg("failed to generate upload URL")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	log.Info().
		Int32("member_id", memberID).
		Str("object_key", objectKey).
		Msg("upload URL generated")

	return c.JSON(http.StatusOK, GenerateUploadURLResponse{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
		ExpiresIn: 900, // 15 minutes
	})
}

// CompleteUploadHandler confirms that an upload is complete and updates the member's image_url
// @Summary Complete profile image upload
// @Description Confirm upload complete and update member's profile image URL
// @Tags upload
// @Accept json
// @Produce json
// @Param request body CompleteUploadRequest true "Complete upload request"
// @Success 200 {object} CompleteUploadResponse "Upload completed"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security SessionAuth
// @Router /upload/profile-image/complete [post]
func (h *UploadHandler) CompleteUploadHandler(c echo.Context) error {
	if !h.s3Service.IsEnabled() {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Storage service not available",
		})
	}

	memberID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	req := new(CompleteUploadRequest)
	if err := helpers.BindAndValidate(c, req); err != nil {
		return err
	}

	log.Info().
		Int32("member_id", memberID).
		Str("object_key", req.ObjectKey).
		Msg("upload completion requested")

	// validate object key format
	if !isValidObjectKey(req.ObjectKey) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid object key"})
	}

	// check if object exists
	exists, err := h.s3Service.ObjectExists(c.Request().Context(), req.ObjectKey)
	if err != nil {
		log.Error().Err(err).Str("object_key", req.ObjectKey).Msg("failed to check object existence")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to verify upload"})
	}
	if !exists {
		log.Warn().
			Int32("member_id", memberID).
			Str("object_key", req.ObjectKey).
			Msg("upload object not found")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Object not found - upload may have failed"})
	}

	// generate download URL
	downloadURL, err := h.s3Service.GenerateDownloadURL(c.Request().Context(), req.ObjectKey)
	if err != nil {
		log.Error().Err(err).Str("object_key", req.ObjectKey).Msg("failed to generate download URL")
		// use public URL as fallback
		downloadURL = h.s3Service.GetPublicURL(req.ObjectKey)
	}

	// update member's image_url
	publicURL := h.s3Service.GetPublicURL(req.ObjectKey)
	err = h.updateMemberImageURL(c.Request().Context(), memberID, publicURL)
	if err != nil {
		log.Error().Err(err).Int32("member_id", memberID).Str("object_key", req.ObjectKey).Msg("failed to update member image URL")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update profile"})
	}

	log.Info().
		Int32("member_id", memberID).
		Str("object_key", req.ObjectKey).
		Str("image_url", publicURL).
		Msg("upload completed")

	return c.JSON(http.StatusOK, CompleteUploadResponse{
		ImageURL:    publicURL,
		DownloadURL: downloadURL,
		Message:     "Profile image updated successfully",
	})
}

// DeleteImageHandler deletes the member's profile image
// @Summary Delete profile image
// @Description Delete the member's profile image from storage
// @Tags upload
// @Produce json
// @Success 200 {object} map[string]string "Image deleted"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security SessionAuth
// @Router /upload/profile-image [delete]
func (h *UploadHandler) DeleteImageHandler(c echo.Context) error {
	if !h.s3Service.IsEnabled() {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Storage service not available",
		})
	}

	memberID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// get current image_url
	ctx := c.Request().Context()
	q := repository.New(h.dbService.GetConnection())

	member, err := q.GetMemberInfoById(ctx, memberID)
	if err != nil {
		log.Error().Err(err).Int32("member_id", memberID).Msg("failed to get member")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get member"})
	}

	if !member.ImageUrl.Valid || member.ImageUrl.String == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No profile image to delete"})
	}

	// extract object key from URL
	objectKey := extractObjectKey(member.ImageUrl.String)
	if objectKey == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid image URL"})
	}

	// delete from S3
	err = h.s3Service.DeleteObject(ctx, objectKey)
	if err != nil {
		log.Error().Err(err).Str("object_key", objectKey).Msg("failed to delete image from S3")
		// continue anyway to clear the database
	}

	// clear image_url in database
	err = h.updateMemberImageURL(ctx, memberID, "")
	if err != nil {
		log.Error().Err(err).Int32("member_id", memberID).Msg("failed to clear member image URL")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to clear profile image"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Profile image deleted successfully",
	})
}

// helper functions

func (h *UploadHandler) updateMemberImageURL(ctx context.Context, memberID int32, imageURL string) error {
	db := h.dbService.GetConnection()

	_, err := db.ExecContext(ctx,
		"UPDATE members SET image_url = ? WHERE id = ?",
		sql.NullString{String: imageURL, Valid: imageURL != ""},
		memberID,
	)
	return err
}

func isValidObjectKey(key string) bool {
	if len(key) < 10 || len(key) > 200 {
		return false
	}
	// must start with profile-images/
	if len(key) < 15 || key[:15] != "profile-images/" {
		return false
	}
	return true
}

func extractObjectKey(url string) string {
	if url == "" {
		return ""
	}
	idx := strings.Index(url, "profile-images/")
	if idx >= 0 {
		return url[idx:]
	}
	return ""
}
