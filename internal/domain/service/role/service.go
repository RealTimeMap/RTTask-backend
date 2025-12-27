package role

import (
	"context"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model/rbac"
	"rttask/internal/domain/repository"

	"go.uber.org/zap"
)

type RoleService struct {
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
	logger   *zap.Logger
}

func NewRoleService(roleRepo repository.RoleRepository, userRepo repository.UserRepository, logger *zap.Logger) *RoleService {
	return &RoleService{roleRepo: roleRepo, userRepo: userRepo, logger: logger}
}

func (s *RoleService) CreateRole(ctx context.Context, input RoleInput) (*rbac.Role, error) {
	// TODO добавить обработку уникальности имени
	// ...
	// ...
	// ...
	// ...

	user, err := s.userRepo.GetUserByIDWithRoles(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	if !user.Can(rbac.RoleCreate) {
		return nil, domainerrors.NewForbiddenError("Dont have permissions to create a role")
	}

	permissions, err := s.validatePermissions(input.Permissions)
	if err != nil {
		return nil, err
	}

	rawRole := &rbac.Role{
		Name:        input.Name,
		Permissions: permissions,
	}

	newRole, err := s.roleRepo.Create(ctx, rawRole)
	if err != nil {
		return nil, err
	}
	return newRole, nil
}

func (s *RoleService) validatePermissions(perms []string) ([]rbac.Permission, error) {
	var permissions []rbac.Permission
	for _, p := range perms {
		if _, exists := rbac.PermissionsRegister[rbac.Permission(p)]; exists {
			permissions = append(permissions, rbac.Permission(p))
		} else {
			return make([]rbac.Permission, 0), domainerrors.NewValidationError("Permission doesn't exist")
		}
	}
	return permissions, nil
}
