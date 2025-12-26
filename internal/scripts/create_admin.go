package scripts

import (
	"context"
	"errors"
	"rttask/internal/config"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model"
	"rttask/internal/domain/repository"
	"rttask/internal/infrastructure/security"

	"go.uber.org/zap"
)

func CreateAdminIfNotExists(ctx context.Context, cfg config.Admin, logger *zap.Logger, repo repository.UserRepository, hasher security.PasswordHasher) {
	existsUser, err := repo.GetUserByEmail(ctx, cfg.Email)
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if errors.As(err, &notFoundErr) || notFoundErr.Type != domainerrors.ErrorTypeNotFound {
		} else {
			panic(err)
		}
	}
	if existsUser != nil {
		logger.Warn("User already exists.", zap.String("email", cfg.Email))
		return
	}

	hashPass, err := hasher.HashPassword(cfg.Password)
	user := &model.User{
		Email:          cfg.Email,
		HashedPassword: hashPass,
		FirstName:      "Admin",
		LastName:       "User",
	}
	adminUser, err := repo.CreateUser(ctx, user)
	if err != nil {
		panic(err)
	}
	logger.Info("admin user created", zap.String("email", adminUser.Email))
}
