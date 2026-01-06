package repository

import (
	"context"
	"rttask/internal/domain/model"
	"rttask/internal/domain/valueobject"
)

type InviteRepository interface {
	Create(ctx context.Context, invite *model.InviteLink) (*model.InviteLink, error)
	GetByToken(ctx context.Context, token string) (*model.InviteLink, error)
	GetAll(ctx context.Context, userID uint, params valueobject.PaginationParams) ([]*model.InviteLink, error)
}
