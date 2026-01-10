package auth

import (
	"context"
	"errors"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model"
	"rttask/internal/domain/repository"
	"rttask/internal/domain/service/file"
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
	fileService     *file.FileService
	passwordHasher  security.PasswordHasher
	jwtManager      security.JWTManager
	accessDuration  time.Duration
	refreshDuration time.Duration
	logger          *zap.Logger
}

func NewAuthService(
	userRepo repository.UserRepository,
	inviteRepo repository.InviteRepository,
	fileService *file.FileService,
	passwordHasher security.PasswordHasher,
	jwtManager security.JWTManager,
	accessDuration time.Duration,
	refreshDuration time.Duration,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		inviteRepo:      inviteRepo,
		fileService:     fileService,
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

func (s *AuthService) Register(ctx context.Context, input RegisterInput, fileInput *file.FileInput) (*model.User, error) {
	if err := s.validateUser(ctx, input.Email.String()); err != nil {
		return nil, err
	}

	invite, err := s.validateInvite(ctx, input.InviteLink) // TODO дальше сделать более сложную проверку
	if err != nil {
		return nil, err
	}

	hashPassword, err := s.hashPassword(input.Password.String())
	if err != nil {
		return nil, err
	}

	var avatar *model.File
	if fileInput != nil {
		uploadAvatar, err := s.fileService.UploadFile(ctx, *fileInput, file.CompanyProfile)
		if err != nil {
			return nil, err
		}
		avatar = uploadAvatar
	}

	newUser := &model.User{
		Email:          input.Email.String(),
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Roles:          invite.Roles,
		HashedPassword: hashPassword,
		Avatar:         avatar,
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

func (s *AuthService) validateInvite(ctx context.Context, inviteLink string) (*model.InviteLink, error) {
	invite, err := s.inviteRepo.GetByToken(ctx, inviteLink) // TODO дальше сделать более сложную проверку
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if errors.As(err, &notFoundErr) && notFoundErr.Type == domainerrors.ErrorTypeNotFound {
			s.logger.Warn("invite token not found",
				zap.String("token", inviteLink),
			)
			return nil, domainerrors.NewNotFoundError("invite", inviteLink)
		}
		// Любая другая ошибка БД
		s.logger.Error("failed to check invite token",
			zap.String("token", inviteLink),
			zap.Error(err),
		)
		return nil, err
	}
	return invite, nil
}

func (s *AuthService) validateUser(ctx context.Context, email string) error {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if !errors.As(err, &notFoundErr) || notFoundErr.Type != domainerrors.ErrorTypeNotFound {
			s.logger.Error("failed to check user existence",
				zap.String("email", email),
				zap.Error(err),
			)
			return err
		}
	}
	if existingUser != nil {
		s.logger.Warn("registration attempt for existing user",
			zap.String("email", email),
		)
		return domainerrors.NewAlreadyExistsError("user", "email", email)
	}
	return nil
}

func (s *AuthService) hashPassword(password string) (string, error) {
	hashPassword, err := s.passwordHasher.HashPassword(password)
	if err != nil {
		s.logger.Error("failed to hash password",
			zap.Error(err),
		)
		return "", domainerrors.NewInternalError("failed to secure password", err)
	}
	return hashPassword, nil
}
