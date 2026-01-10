package invite

import (
	"context"
	"fmt"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model"
	"rttask/internal/domain/model/rbac"
	"rttask/internal/domain/repository"
	"rttask/internal/domain/valueobject"

	"go.uber.org/zap"
)

type InviteService struct {
	inviteRepo repository.InviteRepository
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	logger     *zap.Logger
}

func NewInviteService(inviteRepo repository.InviteRepository, userRepo repository.UserRepository, roleRepo repository.RoleRepository, logger *zap.Logger) *InviteService {
	return &InviteService{
		inviteRepo: inviteRepo,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		logger:     logger,
	}
}

// CreateInvite Создание инвайт ссылки
func (s *InviteService) CreateInvite(ctx context.Context, input InviteInput, userID uint) (*model.InviteLink, error) {
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
func (s *InviteService) GetAllInvites(ctx context.Context, userID uint, params valueobject.PaginationParams) ([]*model.InviteLink, error) {
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

func (s *InviteService) validateUser(ctx context.Context, userID uint, permission rbac.Permission) error {
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

func (s *InviteService) validateRoles(ctx context.Context, input InviteInput) ([]rbac.Role, error) {
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
