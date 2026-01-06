package postgres

import (
	"context"
	"rttask/internal/domain/model"
	"rttask/internal/domain/repository"
	"rttask/internal/domain/valueobject"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PgInviteRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewPgInviteRepository(db *gorm.DB, logger *zap.Logger) repository.InviteRepository {
	return &PgInviteRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PgInviteRepository) Create(ctx context.Context, invite *model.InviteLink) (*model.InviteLink, error) {
	r.logger.Info("start inviteRepository.Create")
	err := r.db.WithContext(ctx).Create(invite).Error
	if err != nil {
		r.logger.Error("Create invite link", zap.Error(err), zap.Any("invite", invite))
		return nil, err
	}
	return invite, nil
}

func (r *PgInviteRepository) GetByToken(ctx context.Context, token string) (*model.InviteLink, error) {
	r.logger.Info("start inviteRepository.GetByToken")
	var invite *model.InviteLink
	err := r.db.WithContext(ctx).Model(&model.InviteLink{}).Preload("Roles").Where("token = ?", token).First(&invite).Error
	if err != nil {
		return nil, MapGormError(err, "invite")
	}
	return invite, nil
}

func (r *PgInviteRepository) GetAll(ctx context.Context, userID uint, params valueobject.PaginationParams) ([]*model.InviteLink, error) {
	var invites []*model.InviteLink
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(params.Offset).
		Limit(params.Limit).
		Find(&invites).Error
	if err != nil {
		return nil, MapGormError(err, "invite")
	}
	return invites, nil
}
