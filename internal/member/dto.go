package member

import (
	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
)

// EmailRequest represents a request with an email parameter
type EmailRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@dlsu.edu.ph"`
}

// IDRequest represents a request with a member ID parameter
type IDRequest struct {
	ID int `json:"id" validate:"required,gt=0" example:"12345678"`
}

// FullInfoMemberResponse represents complete member information
type FullInfoMemberResponse struct {
	ID            int32                  `json:"id" example:"12345678"`
	Email         string                 `json:"email" example:"user@dlsu.edu.ph"`
	FullName      string                 `json:"full_name" example:"Juan Dela Cruz"`
	Nickname      helpers.NullableString `json:"nickname"`
	ImageURL      helpers.NullableString `json:"image_url"`
	CommitteeID   helpers.NullableString `json:"committee_id"`
	CommitteeName helpers.NullableString `json:"committee_name"`
	DivisionID    helpers.NullableString `json:"division_id"`
	DivisionName  helpers.NullableString `json:"division_name"`
	PositionID    helpers.NullableString `json:"position_id"`
	PositionName  helpers.NullableString `json:"position_name"`
	HouseName     helpers.NullableString `json:"house_name"`
	ContactNumber helpers.NullableString `json:"contact_number"`
	College       helpers.NullableString `json:"college"`
	Program       helpers.NullableString `json:"program"`
	Interests     helpers.NullableString `json:"interests"`
	Discord       helpers.NullableString `json:"discord"`
	FbLink        helpers.NullableString `json:"fb_link"`
	Telegram      helpers.NullableString `json:"telegram"`
}

func toFullInfoMemberResponse(m repository.GetMemberInfoRow) FullInfoMemberResponse {
	return FullInfoMemberResponse{
		ID:            m.ID,
		Email:         m.Email,
		FullName:      m.FullName,
		Nickname:      helpers.NullableString{NullString: m.Nickname},
		ImageURL:      helpers.NullableString{NullString: m.ImageUrl},
		CommitteeID:   helpers.NullableString{NullString: m.CommitteeID},
		CommitteeName: helpers.NullableString{NullString: m.CommitteeName},
		DivisionID:    helpers.NullableString{NullString: m.DivisionID},
		DivisionName:  helpers.NullableString{NullString: m.DivisionName},
		PositionID:    helpers.NullableString{NullString: m.PositionID},
		PositionName:  helpers.NullableString{NullString: m.PositionName},
		HouseName:     helpers.NullableString{NullString: m.HouseName},
		ContactNumber: helpers.NullableString{NullString: m.ContactNumber},
		College:       helpers.NullableString{NullString: m.College},
		Program:       helpers.NullableString{NullString: m.Program},
		Interests:     helpers.NullableString{NullString: m.Interests},
		Discord:       helpers.NullableString{NullString: m.Discord},
		FbLink:        helpers.NullableString{NullString: m.FbLink},
		Telegram:      helpers.NullableString{NullString: m.Telegram},
	}
}

// MemberResponse represents basic member information
type MemberResponse struct {
	ID            int32                  `json:"id" example:"12345678"`
	FullName      string                 `json:"full_name" example:"Juan Dela Cruz"`
	Nickname      helpers.NullableString `json:"nickname"`
	Email         string                 `json:"email" example:"user@dlsu.edu.ph"`
	Telegram      helpers.NullableString `json:"telegram"`
	PositionID    helpers.NullableString `json:"position_id"`
	CommitteeID   helpers.NullableString `json:"committee_id"`
	College       helpers.NullableString `json:"college"`
	Program       helpers.NullableString `json:"program"`
	Discord       helpers.NullableString `json:"discord"`
	Interests     helpers.NullableString `json:"interests"`
	ContactNumber helpers.NullableString `json:"contact_number"`
	FbLink        helpers.NullableString `json:"fb_link"`
	ImageURL      helpers.NullableString `json:"image_url"`
	HouseName     helpers.NullableString `json:"house_name"`
}

func toMemberResponse(m repository.ListMembersRow) MemberResponse {
	return MemberResponse{
		ID:            m.ID,
		FullName:      m.FullName,
		Nickname:      helpers.NullableString{NullString: m.Nickname},
		Email:         m.Email,
		Telegram:      helpers.NullableString{NullString: m.Telegram},
		PositionID:    helpers.NullableString{NullString: m.PositionID},
		CommitteeID:   helpers.NullableString{NullString: m.CommitteeID},
		College:       helpers.NullableString{NullString: m.College},
		Program:       helpers.NullableString{NullString: m.Program},
		Discord:       helpers.NullableString{NullString: m.Discord},
		Interests:     helpers.NullableString{NullString: m.Interests},
		ContactNumber: helpers.NullableString{NullString: m.ContactNumber},
		FbLink:        helpers.NullableString{NullString: m.FbLink},
		ImageURL:      helpers.NullableString{NullString: m.ImageUrl},
		HouseName:     helpers.NullableString{NullString: m.HouseName},
	}
}

// UpdateSelfRequest represents a request to update own profile
// Fields: nickname, telegram, discord, interests, contact_number, fb_link, image_url
type UpdateSelfRequest struct {
	Nickname      *string `json:"nickname" validate:"omitempty,max=100"`
	Telegram      *string `json:"telegram" validate:"omitempty,max=100"`
	Discord       *string `json:"discord" validate:"omitempty,max=32"`
	Interests     *string `json:"interests"`
	ContactNumber *string `json:"contact_number" validate:"omitempty,max=32"`
	FbLink        *string `json:"fb_link" validate:"omitempty,max=255"`
	ImageURL      *string `json:"image_url" validate:"omitempty,max=512"`
}

// UpdateMemberRequest represents a request to update another member's profile
// All fields are editable by authorized users
type UpdateMemberRequest struct {
	FullName      *string `json:"full_name" validate:"omitempty,max=255"`
	Nickname      *string `json:"nickname" validate:"omitempty,max=100"`
	Email         *string `json:"email" validate:"omitempty,email"`
	PositionID    *string `json:"position_id" validate:"omitempty,max=10"`
	CommitteeID   *string `json:"committee_id" validate:"omitempty,max=10"`
	College       *string `json:"college" validate:"omitempty,max=255"`
	Program       *string `json:"program" validate:"omitempty,max=255"`
	HouseID       *int    `json:"house_id" validate:"omitempty,gt=0"`
	Telegram      *string `json:"telegram" validate:"omitempty,max=100"`
	Discord       *string `json:"discord" validate:"omitempty,max=32"`
	Interests     *string `json:"interests"`
	ContactNumber *string `json:"contact_number" validate:"omitempty,max=32"`
	FbLink        *string `json:"fb_link" validate:"omitempty,max=255"`
	ImageURL      *string `json:"image_url" validate:"omitempty,max=512"`
}
