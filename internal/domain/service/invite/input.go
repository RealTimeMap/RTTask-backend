package invite

import "github.com/google/uuid"

type InviteInput struct {
	Token    string
	RolesIDs []uint
}

func NewInviteInput(rolesIDs []uint) InviteInput {
	return InviteInput{
		Token:    uuid.New().String(),
		RolesIDs: rolesIDs,
	}
}
