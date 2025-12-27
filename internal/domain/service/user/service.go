package user

import (
	"context"
	"errors"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model"
	"rttask/internal/domain/repository"

	"go.uber.org/zap"
)

type UserService struct {
	userRepo repository.UserRepository
	// TaskRepo
	// RoleRepo
	logger *zap.Logger
}

func (s *UserService) GetUserProfile(ctx context.Context, userID uint) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if errors.As(err, &notFoundErr) && notFoundErr.Type == domainerrors.ErrorTypeNotFound {
			return nil, domainerrors.ErrInvalidCredentials
		}
		return nil, err
	}
	return user, nil
}
