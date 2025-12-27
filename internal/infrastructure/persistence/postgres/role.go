package postgres

import (
	"context"
	"rttask/internal/domain/model/rbac"
	"rttask/internal/domain/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PgRoleRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewPgRoleRepository(db *gorm.DB, logger *zap.Logger) repository.RoleRepository {
	return &PgRoleRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PgRoleRepository) Create(ctx context.Context, role *rbac.Role) (*rbac.Role, error) {
	err := r.db.WithContext(ctx).Create(role).Error
	if err != nil {
		return nil, MapGormError(err, "role")
	}
	return role, nil
}

func (r *PgRoleRepository) Update(ctx context.Context, role *rbac.Role) (*rbac.Role, error) {
	err := r.db.WithContext(ctx).Save(role).Error
	if err != nil {
		return nil, MapGormError(err, "role")
	}
	return role, nil
}

func (r *PgRoleRepository) GetByName(ctx context.Context, name string) (*rbac.Role, error) {
	var role *rbac.Role
	err := r.db.WithContext(ctx).First(&role, "name = ?", name).Error
	if err != nil {
		return nil, MapGormError(err, "role")
	}
	return role, nil
}
