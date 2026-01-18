package invite

import (
	"rttask/internal/domain/model"

	"github.com/google/uuid"
)

type Input struct {
	Token       string
	Description *string
	RolesIDs    []uint
}

func NewInviteInput(description *string, rolesIDs []uint) Input {
	return Input{
		Token:       uuid.New().String(),
		Description: description,
		RolesIDs:    rolesIDs,
	}
}

type UpdateInput struct {
	Description    *string
	RolesIDs       []uint
	DeleteRolesIDs []uint
}

func (i UpdateInput) ApplyTo(invite *model.InviteLink) {
	if i.Description != nil {
		invite.Description = i.Description
	}
}
