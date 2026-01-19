package scripts

import (
	"context"
	"errors"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model/rbac"
	"rttask/internal/domain/repository"

	"go.uber.org/zap"
)

// DefaultRoles содержит определения стандартных системных ролей
var DefaultRoles = []rbac.Role{
	{
		Name:     "Админ",
		IsSystem: true,
		IsActive: true,
		Permissions: []rbac.Permission{
			// Все права - администратор имеет полный доступ
			rbac.InviteCreate,
			rbac.InviteDelete,

			rbac.TaskCreate,
			rbac.TaskDelete,
			rbac.TaskUpdate,
			rbac.TaskAssign,
			rbac.TaskChangeStatus,

			rbac.RoleCreate,
			rbac.RoleUpdate,
			rbac.RoleDelete,
			rbac.RoleAssign,

			rbac.UserUpdate,
			rbac.UserDelete,

			rbac.CompanyCreate,
			rbac.CompanyUpdate,
			rbac.CompanyDelete,

			rbac.CommentCreate,
			rbac.CommentUpdate,
			rbac.CommentDelete,
		},
	},
	{
		Name:     "Менеджер",
		IsSystem: true,
		IsActive: true,
		Permissions: []rbac.Permission{
			// Управление инвайтами
			rbac.InviteCreate,
			rbac.InviteDelete,

			// Полное управление задачами
			rbac.TaskCreate,
			rbac.TaskDelete,
			rbac.TaskUpdate,
			rbac.TaskAssign,
			rbac.TaskChangeStatus,

			// Управление компаниями
			rbac.CompanyCreate,
			rbac.CompanyUpdate,
			rbac.CompanyDelete,

			// Управление комментариями
			rbac.CommentCreate,
			rbac.CommentUpdate,
			rbac.CommentDelete,
		},
	},
	{
		Name:        "Сотрудник",
		IsSystem:    true,
		IsActive:    true,
		Permissions: []rbac.Permission{},
	},
}

func CreateDefaultRolesIfNotExist(ctx context.Context, logger *zap.Logger, roleRepo repository.RoleRepository) {
	for _, role := range DefaultRoles {
		existRole, err := roleRepo.GetByName(ctx, role.Name)
		if err != nil {
			var notFoundError *domainerrors.DomainError
			if errors.As(err, &notFoundError) && notFoundError.Type == domainerrors.ErrorTypeNotFound {
			} else {
				panic(err)
			}
		}
		if existRole != nil {
			logger.Info("role already exist", zap.String("name", role.Name), zap.Any("role", existRole))
			continue
		}
		newRole := &rbac.Role{
			Name:        role.Name,
			IsSystem:    role.IsSystem,
			IsActive:    role.IsActive,
			Permissions: role.Permissions,
		}
		createdRole, err := roleRepo.Create(ctx, newRole)
		if err != nil {
			panic(err)
		}
		logger.Info("role created", zap.String("name", role.Name), zap.Uint("ID", createdRole.ID))
	}
}
