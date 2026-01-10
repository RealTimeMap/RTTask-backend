package invite

import (
	"github.com/google/uuid"
)

type InviteInput struct {
	Token       string
	Description *string
	RolesIDs    []uint
}

func NewInviteInput(description *string, rolesIDs []uint) InviteInput {
	return InviteInput{
		Token:       uuid.New().String(),
		Description: description,
		RolesIDs:    rolesIDs,
	}
}
