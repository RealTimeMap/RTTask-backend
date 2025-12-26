package app

import (
	"rttask/internal/config"
	"rttask/internal/domain/service/auth"
	"rttask/internal/domain/service/invite"
	"rttask/internal/infrastructure/persistence/postgres"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/http/response"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	AuthService   *auth.AuthService
	InviteService *invite.InviteService
	JWTManager    security.JWTManager
	Mapper        *response.ErrorMapper
}

func NewContainer(cfg config.Config, db *gorm.DB, logger *zap.Logger) *Container {
	// Репозитории

	userRepo := postgres.NewPgUserRepository(db, logger)
	inviteRepo := postgres.NewPgInviteRepository(db, logger)
	// JWT хелперы

	passwordHasher := security.NewBcryptHasher()
	manager := security.NewCustomJWTManager(cfg.JWT.Secret)
	mapper := response.NewErrorMapper()

	// Сервисы

	authService := auth.NewAuthService(userRepo, inviteRepo, passwordHasher, manager, cfg.JWT.AccessTokenTimeDuration(), cfg.JWT.RefreshTokenTimeDuration(), logger)
	inviteService := invite.NewInviteService(inviteRepo, userRepo, logger)
	return &Container{
		AuthService:   authService,
		InviteService: inviteService,

		JWTManager: manager,
		Mapper:     mapper,
	}
}
