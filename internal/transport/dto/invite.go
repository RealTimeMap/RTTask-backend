package dto

import (
	"rttask/internal/domain/model"
	"rttask/internal/domain/service/invite"
)

type InviteRequest struct {
	Description *string `json:"description,omitempty"`
	RolesIDs    []uint  `json:"rolesIds"`
}

// RESPONSE

type InviteResponse struct {
	ID          uint     `json:"id"`
	Token       string   `json:"token"`
	Description *string  `json:"description"`
	Roles       []string `json:"roles"`
}

func NewInviteResponse(invite *model.InviteLink) InviteResponse {
	var roles []string
	for _, role := range invite.Roles {
		roles = append(roles, role.Name)
	}
	return InviteResponse{
		ID:          invite.ID,
		Token:       invite.Token,
		Description: invite.Description,
		Roles:       roles,
	}
}

type InviteUpdateRequest struct {
	Description    *string `json:"description"`
	RolesIDs       []uint  `json:"rolesIds"`
	DeleteRolesIDs []uint  `json:"deleteRolesIds"`
}

func (i *InviteUpdateRequest) ToInput() invite.UpdateInput {
	return invite.UpdateInput{
		Description:    i.Description,
		RolesIDs:       i.RolesIDs,
		DeleteRolesIDs: i.DeleteRolesIDs,
	}
}
