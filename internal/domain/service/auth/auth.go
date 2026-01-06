package auth

import (
	"context"
	"errors"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model"
	"rttask/internal/domain/repository"
	"rttask/internal/infrastructure/security"
	"time"

	"go.uber.org/zap"
)

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Type         string `json:"type"`
}

func NewTokens(accessToken, refreshToken string) *Tokens {
	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Type:         "Bearer",
	}
}

type AuthService struct {
	userRepo        repository.UserRepository
	inviteRepo      repository.InviteRepository
	passwordHasher  security.PasswordHasher
	jwtManager      security.JWTManager
	accessDuration  time.Duration
	refreshDuration time.Duration
	logger          *zap.Logger
}

func NewAuthService(
	userRepo repository.UserRepository,
	inviteRepo repository.InviteRepository,
	passwordHasher security.PasswordHasher,
	jwtManager security.JWTManager,
	accessDuration time.Duration,
	refreshDuration time.Duration,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		inviteRepo:      inviteRepo,
		passwordHasher:  passwordHasher,
		jwtManager:      jwtManager,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
		logger:          logger,
	}
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*Tokens, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, input.Email.String())
	if err != nil {
		// Проверяем тип ошибки
		var notFoundErr *domainerrors.DomainError
		if errors.As(err, &notFoundErr) && notFoundErr.Type == domainerrors.ErrorTypeNotFound {
			s.logger.Warn("login attempt for non-existent user",
				zap.String("email", input.Email.String()),
			)
			return nil, domainerrors.ErrInvalidCredentials
		}
		s.logger.Error("failed to get user by email",
			zap.String("email", input.Email.String()),
			zap.Error(err),
		)
		return nil, err
	}

	if err := s.passwordHasher.CheckPassword(user.HashedPassword, input.Password.String()); err != nil {
		s.logger.Warn("invalid password attempt",
			zap.String("email", input.Email.String()),
			zap.Uint("userID", user.ID),
		)
		return nil, domainerrors.ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		s.logger.Error("failed to generate tokens",
			zap.Uint("userID", user.ID),
			zap.Error(err),
		)
		return nil, domainerrors.NewInternalError("failed to generate authentication tokens", err)
	}

	s.logger.Info("user logged in successfully",
		zap.Uint("userID", user.ID),
		zap.String("email", user.Email),
	)

	return tokens, nil
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*model.User, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, input.Email.String())
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if !errors.As(err, &notFoundErr) || notFoundErr.Type != domainerrors.ErrorTypeNotFound {
			s.logger.Error("failed to check user existence",
				zap.String("email", input.Email.String()),
				zap.Error(err),
			)
			return nil, err
		}
	}
	invite, err := s.inviteRepo.GetByToken(ctx, input.InviteLink) // TODO дальше сделать более сложную проверку
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if errors.As(err, &notFoundErr) && notFoundErr.Type == domainerrors.ErrorTypeNotFound {
			s.logger.Warn("invite token not found",
				zap.String("token", input.InviteLink),
			)
			return nil, domainerrors.NewNotFoundError("invite", input.InviteLink)
		}
		// Любая другая ошибка БД
		s.logger.Error("failed to check invite token",
			zap.String("token", input.InviteLink),
			zap.Error(err),
		)
		return nil, err
	}
	if existingUser != nil {
		s.logger.Warn("registration attempt for existing user",
			zap.String("email", input.Email.String()),
		)
		return nil, domainerrors.NewAlreadyExistsError("user", "email", input.Email.String())
	}

	hashPassword, err := s.passwordHasher.HashPassword(input.Password.String())
	if err != nil {
		s.logger.Error("failed to hash password",
			zap.Error(err),
		)
		return nil, domainerrors.NewInternalError("failed to secure password", err)
	}

	newUser := &model.User{
		Email:          input.Email.String(),
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Roles:          invite.Roles,
		HashedPassword: hashPassword,
	}

	createdUser, err := s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		s.logger.Error("failed to create user",
			zap.String("email", input.Email.String()),
			zap.Error(err),
		)
		return nil, err
	}

	// 4. Успех
	s.logger.Info("user registered successfully",
		zap.Uint("userID", createdUser.ID),
		zap.String("email", createdUser.Email),
	)

	return createdUser, nil
}

func (s *AuthService) generateTokens(user *model.User) (*Tokens, error) {
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, security.AccessToken, s.accessDuration)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtManager.GenerateToken(user.ID, user.Email, security.RefreshToken, s.refreshDuration)
	if err != nil {
		return nil, err
	}
	return NewTokens(accessToken, refreshToken), nil
}
