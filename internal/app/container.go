package app

import (
	"rttask/internal/config"
	"rttask/internal/domain/repository"
	"rttask/internal/domain/service/auth"
	"rttask/internal/domain/service/invite"
	"rttask/internal/domain/service/role"
	"rttask/internal/infrastructure/persistence/postgres"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/http/response"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	AuthService   *auth.AuthService
	InviteService *invite.InviteService
	RoleService   *role.RoleService

	JWTManager security.JWTManager
	Mapper     *response.ErrorMapper
	Hasher     security.PasswordHasher

	UserRepository repository.UserRepository
}

func NewContainer(cfg config.Config, db *gorm.DB, logger *zap.Logger) *Container {
	// Репозитории

	userRepo := postgres.NewPgUserRepository(db, logger)
	inviteRepo := postgres.NewPgInviteRepository(db, logger)
	roleRepo := postgres.NewPgRoleRepository(db, logger)
	// JWT хелперы

	passwordHasher := security.NewBcryptHasher()
	manager := security.NewCustomJWTManager(cfg.JWT.Secret)
	mapper := response.NewErrorMapper()

	// Сервисы

	authService := auth.NewAuthService(userRepo, inviteRepo, passwordHasher, manager, cfg.JWT.AccessTokenTimeDuration(), cfg.JWT.RefreshTokenTimeDuration(), logger)
	inviteService := invite.NewInviteService(inviteRepo, userRepo, logger)
	roleService := role.NewRoleService(roleRepo, userRepo, logger)

	return &Container{
		AuthService:   authService,
		InviteService: inviteService,
		RoleService:   roleService,

		JWTManager: manager,
		Mapper:     mapper,

		UserRepository: userRepo,
		Hasher:         passwordHasher,
	}
}
