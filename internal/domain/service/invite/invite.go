package invite

import (
	"context"
	"fmt"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model"
	"rttask/internal/domain/model/rbac"
	"rttask/internal/domain/repository"
	"rttask/internal/domain/valueobject"
	"slices"

	"go.uber.org/zap"
)

type Service struct {
	inviteRepo repository.InviteRepository
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	logger     *zap.Logger
}

func NewInviteService(inviteRepo repository.InviteRepository, userRepo repository.UserRepository, roleRepo repository.RoleRepository, logger *zap.Logger) *Service {
	return &Service{
		inviteRepo: inviteRepo,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		logger:     logger,
	}
}

// CreateInvite Создание инвайт ссылки
func (s *Service) CreateInvite(ctx context.Context, input Input, userID uint) (*model.InviteLink, error) {
	s.logger.Info("start InviteService.CreateInvite", zap.Any("input", input))

	// Проверка пользоваеля
	err := s.validateUser(ctx, userID, rbac.RoleCreate)
	if err != nil {
		return nil, err
	}

	// Виладиция ролей
	roles, err := s.validateRoles(ctx, input)
	if err != nil {
		return nil, err
	}

	// Создание записи
	invite := &model.InviteLink{Token: input.Token, Description: input.Description, Roles: roles, UserID: userID}
	invite, err = s.inviteRepo.Create(ctx, invite)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// GetAllInvites Получение инвайт ссылок с пагинацией
func (s *Service) GetAllInvites(ctx context.Context, userID uint, params valueobject.PaginationParams) ([]*model.InviteLink, error) {
	user, err := s.userRepo.GetUserByIDWithRoles(ctx, userID)
	if err != nil {
		return nil, domainerrors.NewUnauthorizedError("Not authorized")
	}
	can := user.Can(rbac.InviteList)
	if !can {
		return nil, domainerrors.NewUnauthorizedError("Not allowed")
	}
	invites, err := s.inviteRepo.GetAll(ctx, userID, params)
	s.logger.Info("invites", zap.Int("count", len(invites)))
	if err != nil {
		return nil, err
	}
	return invites, nil
}

func (s *Service) UpdateInvite(ctx context.Context, input UpdateInput, inviteID, userID uint) (*model.InviteLink, error) {
	s.logger.Info("start InviteService.UpdateInvite", zap.Any("input", input))

	if err := s.validateUser(ctx, userID, rbac.RoleUpdate); err != nil {
		return nil, err
	}

	invite, err := s.inviteRepo.GetByID(ctx, inviteID)
	if err != nil {
		return nil, err
	}

	// Применяем изменения полей
	input.ApplyTo(invite)

	// Удаляем роли
	if len(input.DeleteRolesIDs) > 0 {
		invite.Roles = s.removeRoles(invite.Roles, input.DeleteRolesIDs)
	}

	// Добавляем новые роли
	if len(input.RolesIDs) > 0 {
		newRoles, err := s.roleRepo.GetByIDs(ctx, input.RolesIDs)
		if err != nil {
			return nil, err
		}
		invite.Roles = append(invite.Roles, newRoles...)
	}

	return s.inviteRepo.Update(ctx, invite)
}

func (s *Service) removeRoles(roles []rbac.Role, deleteIDs []uint) []rbac.Role {
	result := make([]rbac.Role, 0, len(roles))
	for _, role := range roles {
		if slices.Contains(deleteIDs, role.ID) {
			continue
		}
		result = append(result, role)
	}
	return result
}

func (s *Service) validateUser(ctx context.Context, userID uint, permission rbac.Permission) error {
	user, err := s.userRepo.GetUserByIDWithRoles(ctx, userID)
	if err != nil {
		return domainerrors.NewUnauthorizedError("Not authorized")
	}
	can := user.CanAll(permission)
	if !can {
		return domainerrors.NewUnauthorizedError("Not allowed") // Todo обработку прав доступа
	}
	return nil
}

func (s *Service) validateRoles(ctx context.Context, input Input) ([]rbac.Role, error) {
	if len(input.RolesIDs) > 0 {
		roles, err := s.roleRepo.GetByIDs(ctx, input.RolesIDs)
		if err != nil {
			return nil, err
		}
		if len(roles) != len(input.RolesIDs) {
			return nil, domainerrors.NewValidationError("Invalid role IDs")
		}
		for _, role := range roles {
			if !role.IsActive {
				return nil, domainerrors.NewValidationError(fmt.Sprintf("Role %s is not active", role.Name))
			}
		}
		return roles, nil
	}
	return nil, domainerrors.NewValidationError("Invalid role IDs")
}
