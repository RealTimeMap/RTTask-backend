package repository

import (
	"context"
	"rttask/internal/domain/model"
)

type InviteRepository interface {
	Create(ctx context.Context, invite *model.InviteLink) (*model.InviteLink, error)
	GetByToken(ctx context.Context, token string) (*model.InviteLink, error)
}
