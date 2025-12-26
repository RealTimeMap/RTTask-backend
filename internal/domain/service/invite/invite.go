package invite

import (
	"context"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model"
	"rttask/internal/domain/model/rbac"
	"rttask/internal/domain/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InviteService struct {
	inviteRepo repository.InviteRepository
	userRepo   repository.UserRepository
	logger     *zap.Logger
}

func NewInviteService(inviteRepo repository.InviteRepository, userRepo repository.UserRepository, logger *zap.Logger) *InviteService {
	return &InviteService{
		inviteRepo: inviteRepo,
		userRepo:   userRepo,
		logger:     logger,
	}
}

func (s *InviteService) CreateInvite(ctx context.Context, userID uint) (*model.InviteLink, error) {
	user, err := s.userRepo.GetUserByIDWithRoles(ctx, userID)
	if err != nil {
		return nil, domainerrors.NewUnauthorizedError("Not authorized")
	}
	can := user.Can(rbac.InviteCreate)
	if !can {
		return nil, domainerrors.NewUnauthorizedError("Not allowed") // Todo обработку прав доступа
	}
	invite := &model.InviteLink{Token: uuid.New().String()}
	invite, err = s.inviteRepo.Create(ctx, invite)
	if err != nil {
		return nil, err
	}
	return invite, nil
}
