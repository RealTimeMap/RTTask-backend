package repository

import (
	"context"
	"rttask/internal/domain/model/rbac"
)

type RoleRepository interface {
	Create(ctx context.Context, role *rbac.Role) (*rbac.Role, error)
	Update(ctx context.Context, role *rbac.Role) (*rbac.Role, error)
}
