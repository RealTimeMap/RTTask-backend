package role

import (
	"context"
	"errors"
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

// CreateRole создание роли если данные прошли валидацию
func (s *RoleService) CreateRole(ctx context.Context, input RoleInput) (*rbac.Role, error) {
	permissions, err := s.validate(ctx, input)
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

// validate Объеденяет все валидацию в 1 функцию
func (s *RoleService) validate(ctx context.Context, input RoleInput) ([]rbac.Permission, error) {
	err := s.checkExistRole(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	err = s.checkUserPermission(ctx, input.UserID, rbac.RoleCreate)
	if err != nil {
		return nil, err
	}
	permissions, err := s.validatePermissions(input.Permissions)
	if err != nil {
		return nil, err
	}
	return permissions, err
}

// validatePermissions Валидирует права и возвращает массив прав
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

// checkExistRole проверяет что создаваемая роль не существует
func (s *RoleService) checkExistRole(ctx context.Context, name string) error {
	role, err := s.roleRepo.GetByName(ctx, name)
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if !errors.As(err, &notFoundErr) && notFoundErr.Type == domainerrors.ErrorTypeNotFound {
			return err
		}
	}
	if role != nil {
		return domainerrors.NewAlreadyExistsError("role", "name", name)
	}
	return nil
}

// checkUserPermission Валидирует права доступа пользователя. может ли он сделать это дейтсвие
func (s *RoleService) checkUserPermission(ctx context.Context, userID uint, permission rbac.Permission) error {
	user, err := s.userRepo.GetUserByIDWithRoles(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return domainerrors.NewUnauthorizedError("Unauthorized")
	}
	if !user.Can(permission) {
		return domainerrors.NewForbiddenError("Dont have permissions to create a role")
	}
	return nil
}
