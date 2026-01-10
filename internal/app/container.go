package app

import (
	"rttask/internal/config"
	"rttask/internal/domain/repository"
	"rttask/internal/domain/service/auth"
	"rttask/internal/domain/service/company"
	"rttask/internal/domain/service/file"
	"rttask/internal/domain/service/invite"
	"rttask/internal/domain/service/role"
	"rttask/internal/domain/service/task"
	"rttask/internal/infrastructure/persistence/postgres"
	"rttask/internal/infrastructure/security"
	"rttask/internal/infrastructure/storage"
	"rttask/internal/transport/http/response"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	AuthService    *auth.AuthService
	InviteService  *invite.InviteService
	RoleService    *role.RoleService
	CompanyService *company.CompanyService
	TaskService    *task.TaskService

	JWTManager security.JWTManager
	Mapper     *response.ErrorMapper
	Hasher     security.PasswordHasher

	UserRepository repository.UserRepository
	RoleRepository repository.RoleRepository
}

func NewContainer(cfg config.Config, db *gorm.DB, logger *zap.Logger) *Container {
	// Репозитории

	userRepo := postgres.NewPgUserRepository(db, logger)
	inviteRepo := postgres.NewPgInviteRepository(db, logger)
	roleRepo := postgres.NewPgRoleRepository(db, logger)
	companyRepo := postgres.NewPgCompanyRepository(db, logger)
	taskRepo := postgres.NewPgTaskRepository(db, logger)
	// JWT хелперы

	passwordHasher := security.NewBcryptHasher()
	manager := security.NewCustomJWTManager(cfg.JWT.Secret)
	mapper := response.NewErrorMapper()

	store := storage.NewLocalStorage("./store")

	// Сервисы
	fileService := file.NewFileService(store, logger)
	authService := auth.NewAuthService(userRepo, inviteRepo, fileService, passwordHasher, manager, cfg.JWT.AccessTokenTimeDuration(), cfg.JWT.RefreshTokenTimeDuration(), logger)
	inviteService := invite.NewInviteService(inviteRepo, userRepo, roleRepo, logger)
	roleService := role.NewRoleService(roleRepo, userRepo, logger)
	companyService := company.NewCompanyService(companyRepo, userRepo, fileService, logger)
	taskService := task.NewTaskService(taskRepo, userRepo, companyRepo, fileService, logger)
	return &Container{
		AuthService:    authService,
		InviteService:  inviteService,
		RoleService:    roleService,
		CompanyService: companyService,
		TaskService:    taskService,

		JWTManager: manager,
		Mapper:     mapper,
		Hasher:     passwordHasher,

		UserRepository: userRepo,
		RoleRepository: roleRepo,
	}
}
