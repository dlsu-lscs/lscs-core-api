package member

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetMemberInfo(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	req := new(EmailRequest)

	if err := c.Bind(req); err != nil {
		slog.Error("Failed to parse request body", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	memberInfo, err := q.GetMemberInfo(ctx, req.Email)
	if err != nil {
		slog.Error("email is not an LSCS member", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Email is not an LSCS member"})
	}

	response := toFullInfoMemberResponse(memberInfo)

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetMemberInfoByID(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	req := new(IdRequest)

	if err := c.Bind(req); err != nil {
		slog.Error("Failed to parse request body", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	memberInfo, err := q.GetMemberInfoById(ctx, int32(req.Id))
	if err != nil {
		slog.Error("id is not an LSCS member", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is not an LSCS member"})
	}

	response := toFullInfoMemberResponse(repository.GetMemberInfoRow(memberInfo))

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetAllMembersHandler(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	queries := repository.New(dbconn)

	members, err := queries.ListMembers(ctx)
	if err != nil {
		slog.Error("Failed to list members", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list members"})
	}

	response := make([]MemberResponse, 0, len(members))
	for _, m := range members {
		response = append(response, toMemberResponse(m))
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) CheckEmailHandler(c echo.Context) error {
	var req EmailRequest

	if err := c.Bind(&req); err != nil {
		slog.Error("invalid request body")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Email == "" {
		slog.Error("email is required")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email is required"})
	}

	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	queries := repository.New(dbconn)
	memberEmail, err := queries.CheckEmailIfMember(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": req.Email,
			}
			return c.JSON(http.StatusNotFound, response)
		}
		slog.Error("Error checking email", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := map[string]any{
		"success": "Email is an LSCS member",
		"state":   "present",
		"email":   memberEmail,
	}
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) CheckIDIfMember(c echo.Context) error {
	var req IdRequest

	if err := c.Bind(&req); err != nil {
		slog.Error("invalid request body")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)
	id, err := q.CheckIdIfMember(c.Request().Context(), int32(req.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]any{
				"error": "Not an LSCS member",
				"state": "absent",
				"id":    req.Id,
			}
			return c.JSON(http.StatusNotFound, response)
		}
		slog.Error("invalid ID", "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "invalid ID"})
	}

	response := map[string]any{
		"success": "ID is an LSCS member",
		"state":   "present",
		"id":      id,
	}
	return c.JSON(http.StatusOK, response)
}