package postgres

import (
	"context"
	"rttask/internal/domain/model"
	"rttask/internal/domain/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PgUserRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewPgUserRepository(db *gorm.DB, logger *zap.Logger) repository.UserRepository {
	return &PgUserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PgUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, MapGormError(err, "user")
	}
	return user, nil
}

func (r *PgUserRepository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user *model.User
	err := r.db.WithContext(ctx).Preload("Roles").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, MapGormError(err, "user")
	}
	return user, nil
}

func (r *PgUserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, MapGormError(err, "user")
	}
	return user, nil
}

func (r *PgUserRepository) GetUserByIDWithRoles(ctx context.Context, id uint) (*model.User, error) {
	var user *model.User
	err := r.db.WithContext(ctx).Preload("Roles").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, MapGormError(err, "user")
	}
	return user, nil
}
