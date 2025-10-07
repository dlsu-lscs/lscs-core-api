package member

import (
	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
)

type EmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type IdRequest struct {
	Id int `json:"id" validate:"required,id"`
}

type FullInfoMemberResponse struct {
	ID            int32                  `json:"id"`
	Email         string                 `json:"email"`
	FullName      string                 `json:"full_name"`
	Nickname      helpers.NullableString `json:"nickname"`
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

type MemberResponse struct {
	ID            int32                  `json:"id"`
	FullName      string                 `json:"full_name"`
	Nickname      helpers.NullableString `json:"nickname"`
	Email         string                 `json:"email"`
	Telegram      helpers.NullableString `json:"telegram"`
	PositionID    helpers.NullableString `json:"position_id"`
	CommitteeID   helpers.NullableString `json:"committee_id"`
	College       helpers.NullableString `json:"college"`
	Program       helpers.NullableString `json:"program"`
	Discord       helpers.NullableString `json:"discord"`
	Interests     helpers.NullableString `json:"interests"`
	ContactNumber helpers.NullableString `json:"contact_number"`
	FbLink        helpers.NullableString `json:"fb_link"`
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
		HouseName:     helpers.NullableString{NullString: m.HouseName},
	}
}
