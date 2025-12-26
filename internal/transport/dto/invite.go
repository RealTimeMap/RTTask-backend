package dto

import "rttask/internal/domain/model"

type InviteRequest struct {
}

// RESPONSE

type InviteResponse struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

func NewInviteResponse(invite *model.InviteLink) InviteResponse {
	return InviteResponse{
		ID:    invite.ID,
		Token: invite.Token,
	}
}
