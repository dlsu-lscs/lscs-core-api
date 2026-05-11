package member

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
)

// GetMemberInfo retrieves detailed member information by email
// @Summary Get member info by email
// @Description Get complete member information using their email address
// @Tags members
// @Accept json
// @Produce json
// @Param request body EmailRequest true "Email Request"
// @Success 200 {object} FullInfoMemberResponse "Member information"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 404 {object} helpers.ErrorResponse "Member not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /member [post]
func (h *Handler) GetMemberInfo(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	req := new(EmailRequest)

	if err := helpers.BindAndValidate(c, req); err != nil {
		return err
	}

	memberInfo, err := q.GetMemberInfo(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": req.Email,
			}
			return c.JSON(http.StatusNotFound, response)
		}
		log.Error().Err(err).Msg("error checking email")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := toFullInfoMemberResponse(memberInfo)

	return c.JSON(http.StatusOK, response)
}

// GetMemberInfoByID retrieves detailed member information by ID
// @Summary Get member info by ID
// @Description Get complete member information using their student ID
// @Tags members
// @Accept json
// @Produce json
// @Param request body IDRequest true "ID Request"
// @Success 200 {object} FullInfoMemberResponse "Member information"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 404 {object} helpers.ErrorResponse "Member not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /member-id [post]
func (h *Handler) GetMemberInfoByID(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	req := new(IDRequest)

	if err := helpers.BindAndValidate(c, req); err != nil {
		return err
	}

	memberInfo, err := q.GetMemberInfoById(ctx, int32(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Int("id", req.ID).Msg("id is not an LSCS member")
			response := map[string]string{
				"error": "Not an LSCS member",
				"state": "absent",
				"id":    strconv.Itoa(req.ID),
			}
			return c.JSON(http.StatusNotFound, response)
		}
		log.Error().Err(err).Msg("error checking id")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := toFullInfoMemberResponse(repository.GetMemberInfoRow(memberInfo))

	return c.JSON(http.StatusOK, response)
}

// GetAllMembersHandler lists all members
// @Summary List all members
// @Description Get a list of all LSCS members with basic information
// @Tags members
// @Produce json
// @Param position query string false "Filter by position IDs (comma-separated or repeated)"
// @Param committee query string false "Filter by committee IDs (comma-separated or repeated)"
// @Success 200 {array} MemberResponse "List of members"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /members [get]
func (h *Handler) GetAllMembersHandler(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	queries := repository.New(dbconn)

	positionValues := append([]string{}, c.QueryParams()["position"]...)
	positionValues = append(positionValues, c.QueryParams()["position_id"]...)
	committeeValues := append([]string{}, c.QueryParams()["committee"]...)
	committeeValues = append(committeeValues, c.QueryParams()["committee_id"]...)

	positionIDs := parseQueryListParam(positionValues)
	committeeIDs := parseQueryListParam(committeeValues)

	members, err := queries.ListMembers(ctx, repository.ListMembersParams{
		PositionIds:  positionIDs,
		CommitteeIds: committeeIDs,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to list members")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list members"})
	}

	response := make([]MemberResponse, 0, len(members))
	for _, m := range members {
		response = append(response, toMemberResponse(m))
	}

	return c.JSON(http.StatusOK, response)
}

func parseQueryListParam(values []string) string {
	if len(values) == 0 {
		return ""
	}

	items := make([]string, 0, len(values))
	for _, value := range values {
		for part := range strings.SplitSeq(value, ",") {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}
			items = append(items, trimmed)
		}
	}

	return strings.Join(items, ",")
}

// CheckEmailHandler checks if an email belongs to an LSCS member
// @Summary Check email membership
// @Description Verify if an email address belongs to an LSCS member
// @Tags members
// @Accept json
// @Produce json
// @Param request body EmailRequest true "Email Request"
// @Success 200 {object} map[string]interface{} "Member exists"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 404 {object} helpers.ErrorResponse "Member not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /check-email [post]
func (h *Handler) CheckEmailHandler(c echo.Context) error {
	var req EmailRequest

	if err := helpers.BindAndValidate(c, &req); err != nil {
		return err
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
		log.Error().Err(err).Msg("error checking email")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := map[string]any{
		"success": "Email is an LSCS member",
		"state":   "present",
		"email":   memberEmail,
	}
	return c.JSON(http.StatusOK, response)
}

// CheckIDIfMember checks if an ID belongs to an LSCS member
// @Summary Check ID membership
// @Description Verify if a student ID belongs to an LSCS member
// @Tags members
// @Accept json
// @Produce json
// @Param request body IDRequest true "ID Request"
// @Success 200 {object} map[string]interface{} "Member exists"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 404 {object} helpers.ErrorResponse "Member not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /check-id [post]
func (h *Handler) CheckIDIfMember(c echo.Context) error {
	var req IDRequest

	if err := helpers.BindAndValidate(c, &req); err != nil {
		return err
	}

	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)
	id, err := q.CheckIdIfMember(c.Request().Context(), int32(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]any{
				"error": "Not an LSCS member",
				"state": "absent",
				"id":    req.ID,
			}
			return c.JSON(http.StatusNotFound, response)
		}
		log.Error().Err(err).Msg("invalid ID")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "invalid ID"})
	}

	response := map[string]any{
		"success": "ID is an LSCS member",
		"state":   "present",
		"id":      id,
	}
	return c.JSON(http.StatusOK, response)
}

// GetMeHandler retrieves the authenticated member's profile
// @Summary Get my profile
// @Description Get the authenticated member's complete profile information
// @Tags members
// @Produce json
// @Success 200 {object} FullInfoMemberResponse "Member profile"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security SessionAuth
// @Router /me [get]
func (h *Handler) GetMeHandler(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	email, ok := c.Get("user_email").(string)
	if !ok || email == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	memberInfo, err := q.GetMemberInfo(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Member not found"})
		}
		log.Error().Err(err).Msg("error getting member profile")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := toFullInfoMemberResponse(repository.GetMemberInfoRow(memberInfo))
	h.resolveImageURL(ctx, &response)

	return c.JSON(http.StatusOK, response)
}

// UpdateMeHandler updates the authenticated member's own profile
// @Summary Update my profile
// @Description Update the authenticated member's profile (self-editable fields only)
// @Tags members
// @Accept json
// @Produce json
// @Param request body UpdateSelfRequest true "Update request"
// @Success 200 {object} FullInfoMemberResponse "Updated profile"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security SessionAuth
// @Router /me [put]
func (h *Handler) UpdateMeHandler(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	email, ok := c.Get("user_email").(string)
	if !ok || email == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// get member ID first
	member, err := q.GetMemberInfo(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Member not found"})
		}
		log.Error().Err(err).Msg("error getting member")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	memberID := member.ID

	// bind and validate request
	req := new(UpdateSelfRequest)
	if err := helpers.BindAndValidate(c, req); err != nil {
		return err
	}

	// prepare update values
	nickname := sql.NullString{Valid: false}
	if req.Nickname != nil {
		nickname.String = *req.Nickname
		nickname.Valid = true
	}

	telegram := sql.NullString{Valid: false}
	if req.Telegram != nil {
		telegram.String = *req.Telegram
		telegram.Valid = true
	}

	discord := sql.NullString{Valid: false}
	if req.Discord != nil {
		discord.String = *req.Discord
		discord.Valid = true
	}

	interests := sql.NullString{Valid: false}
	if req.Interests != nil {
		interests.String = *req.Interests
		interests.Valid = true
	}

	contactNumber := sql.NullString{Valid: false}
	if req.ContactNumber != nil {
		contactNumber.String = *req.ContactNumber
		contactNumber.Valid = true
	}

	fbLink := sql.NullString{Valid: false}
	if req.FbLink != nil {
		fbLink.String = *req.FbLink
		fbLink.Valid = true
	}

	imageURL := sql.NullString{Valid: false}
	if req.ImageURL != nil {
		imageURL.String = *req.ImageURL
		imageURL.Valid = true
	}

	// execute update
	err = q.UpdateMemberSelf(ctx, repository.UpdateMemberSelfParams{
		Nickname:      nickname,
		Telegram:      telegram,
		Discord:       discord,
		Interests:     interests,
		ContactNumber: contactNumber,
		FbLink:        fbLink,
		ImageUrl:      imageURL,
		ID:            memberID,
	})
	if err != nil {
		log.Error().Err(err).Int32("member_id", memberID).Msg("error updating member")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update profile"})
	}

	// fetch updated profile
	updatedMember, err := q.GetMemberInfoById(ctx, memberID)
	if err != nil {
		log.Error().Err(err).Int32("member_id", memberID).Msg("error fetching updated member")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch updated profile"})
	}

	response := toFullInfoMemberResponse(repository.GetMemberInfoRow(updatedMember))
	h.resolveImageURL(ctx, &response)

	return c.JSON(http.StatusOK, response)
}

// GetMemberByIDHandler retrieves a member's profile by ID
// @Summary Get member by ID
// @Description Get a member's complete profile information by their ID
// @Tags members
// @Produce json
// @Param id path int true "Member ID"
// @Success 200 {object} FullInfoMemberResponse "Member profile"
// @Failure 400 {object} helpers.ErrorResponse "Invalid ID"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 404 {object} helpers.ErrorResponse "Member not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security SessionAuth
// @Router /members/{id} [get]
func (h *Handler) GetMemberByIDHandler(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid member ID"})
	}

	memberInfo, err := q.GetMemberInfoById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Member not found"})
		}
		log.Error().Err(err).Int64("id", id).Msg("error getting member")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := toFullInfoMemberResponse(repository.GetMemberInfoRow(memberInfo))
	h.resolveImageURL(ctx, &response)

	return c.JSON(http.StatusOK, response)
}

// UpdateMemberByIDHandler updates a member's profile (admin/authorized only)
// @Summary Update member by ID
// @Description Update a member's profile (requires authorization)
// @Tags members
// @Accept json
// @Produce json
// @Param id path int true "Member ID"
// @Param request body UpdateMemberRequest true "Update request"
// @Success 200 {object} FullInfoMemberResponse "Updated profile"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 403 {object} helpers.ErrorResponse "Forbidden - cannot edit this member"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security SessionAuth
// @Router /members/{id} [put]
func (h *Handler) UpdateMemberByIDHandler(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	// get actor ID from context
	actorID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// get target member ID from path
	idStr := c.Param("id")
	targetID, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid member ID"})
	}

	// check if target member exists
	_, err = q.GetMemberInfoById(ctx, int32(targetID))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Member not found"})
		}
		log.Error().Err(err).Int64("id", targetID).Msg("error getting member")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	// bind and validate request
	req := new(UpdateMemberRequest)
	if err := helpers.BindAndValidate(c, req); err != nil {
		return err
	}

	// prepare update values (only set non-nil fields)
	// FullName and Email are required strings, use existing value if not provided
	fullName := ""
	if req.FullName != nil {
		fullName = *req.FullName
	} else {
		// fetch existing full name
		existing, _ := q.GetMemberInfoById(ctx, int32(targetID))
		fullName = existing.FullName
	}

	email := ""
	if req.Email != nil {
		email = *req.Email
	} else {
		// fetch existing email
		existing, _ := q.GetMemberInfoById(ctx, int32(targetID))
		email = existing.Email
	}

	nickname := sql.NullString{Valid: false}
	if req.Nickname != nil {
		nickname.String = *req.Nickname
		nickname.Valid = true
	}

	positionID := sql.NullString{Valid: false}
	if req.PositionID != nil {
		positionID.String = *req.PositionID
		positionID.Valid = true
	}

	committeeID := sql.NullString{Valid: false}
	if req.CommitteeID != nil {
		committeeID.String = *req.CommitteeID
		committeeID.Valid = true
	}

	college := sql.NullString{Valid: false}
	if req.College != nil {
		college.String = *req.College
		college.Valid = true
	}

	program := sql.NullString{Valid: false}
	if req.Program != nil {
		program.String = *req.Program
		program.Valid = true
	}

	houseID := sql.NullInt32{Valid: false}
	if req.HouseID != nil {
		houseID.Int32 = int32(*req.HouseID)
		houseID.Valid = true
	}

	telegram := sql.NullString{Valid: false}
	if req.Telegram != nil {
		telegram.String = *req.Telegram
		telegram.Valid = true
	}

	discord := sql.NullString{Valid: false}
	if req.Discord != nil {
		discord.String = *req.Discord
		discord.Valid = true
	}

	interests := sql.NullString{Valid: false}
	if req.Interests != nil {
		interests.String = *req.Interests
		interests.Valid = true
	}

	contactNumber := sql.NullString{Valid: false}
	if req.ContactNumber != nil {
		contactNumber.String = *req.ContactNumber
		contactNumber.Valid = true
	}

	fbLink := sql.NullString{Valid: false}
	if req.FbLink != nil {
		fbLink.String = *req.FbLink
		fbLink.Valid = true
	}

	imageURL := sql.NullString{Valid: false}
	if req.ImageURL != nil {
		imageURL.String = *req.ImageURL
		imageURL.Valid = true
	}

	// execute update
	err = q.UpdateMemberById(ctx, repository.UpdateMemberByIdParams{
		FullName:      fullName,
		Nickname:      nickname,
		Email:         email,
		PositionID:    positionID,
		CommitteeID:   committeeID,
		College:       college,
		Program:       program,
		HouseID:       houseID,
		Telegram:      telegram,
		Discord:       discord,
		Interests:     interests,
		ContactNumber: contactNumber,
		FbLink:        fbLink,
		ImageUrl:      imageURL,
		ID:            int32(targetID),
	})
	if err != nil {
		log.Error().Err(err).Int32("actor_id", actorID).Int64("target_id", targetID).Msg("error updating member")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update member"})
	}

	// fetch updated profile
	updatedMember, err := q.GetMemberInfoById(ctx, int32(targetID))
	if err != nil {
		log.Error().Err(err).Int64("id", targetID).Msg("error fetching updated member")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch updated profile"})
	}

	response := toFullInfoMemberResponse(repository.GetMemberInfoRow(updatedMember))
	h.resolveImageURL(ctx, &response)

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) resolveImageURL(ctx context.Context, response *FullInfoMemberResponse) {
	if h.s3Service == nil || !h.s3Service.IsEnabled() {
		return
	}
	if !response.ImageURL.Valid || response.ImageURL.String == "" {
		return
	}
	resolved := h.s3Service.ResolveImageURL(ctx, response.ImageURL.String)
	if resolved != response.ImageURL.String {
		response.ImageURL = helpers.NullableString{NullString: sql.NullString{String: resolved, Valid: true}}
	}
}
