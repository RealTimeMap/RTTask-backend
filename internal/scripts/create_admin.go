package scripts

import (
	"context"
	"errors"
	"rttask/internal/config"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model"
	"rttask/internal/domain/model/rbac"
	"rttask/internal/domain/repository"
	"rttask/internal/infrastructure/security"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateAdminIfNotExists(ctx context.Context, cfg config.Admin, logger *zap.Logger, repo repository.UserRepository, hasher security.PasswordHasher) {
	existsUser, err := repo.GetUserByEmail(ctx, cfg.Email)
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if errors.As(err, &notFoundErr) && notFoundErr.Type == domainerrors.ErrorTypeNotFound {
		} else {
			panic(err)
		}
	}
	if existsUser != nil {
		logger.Warn("Admin user already exists.", zap.String("email", cfg.Email))
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

func CreateAdminRoleIfNotExists(ctx context.Context, logger *zap.Logger, roleRepo repository.RoleRepository) {
	existsRole, err := roleRepo.GetByName(ctx, "admin")
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if errors.As(err, &notFoundErr) && notFoundErr.Type == domainerrors.ErrorTypeNotFound {
		} else {
			panic(err)
		}
	}
	if existsRole != nil {
		logger.Warn("Admin role already exists.", zap.String("role", "admin"))
		return
	}

	var permissions []rbac.Permission
	for _, p := range rbac.PermissionsRegister {
		permissions = append(permissions, rbac.Permission(p.Name))
	}

	// Создаем системную роль админа
	adminRole := &rbac.Role{
		Name:        "admin",
		Permissions: permissions,
		IsSystem:    true,
		IsActive:    true,
	}

	createdRole, err := roleRepo.Create(ctx, adminRole)
	if err != nil {
		panic(err)
	}
	logger.Info("admin role created", zap.String("role", createdRole.Name), zap.Int("permissions_count", len(createdRole.Permissions)))
}

func AssignAdminRoleToAdmin(ctx context.Context, cfg config.Admin, logger *zap.Logger, roleRepo repository.RoleRepository, db *gorm.DB) {
	// Получаем роль админа
	adminRole, err := roleRepo.GetByName(ctx, "admin")
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if errors.As(err, &notFoundErr) && notFoundErr.Type == domainerrors.ErrorTypeNotFound {
			logger.Warn("Admin role not found, skipping role assignment", zap.String("role", "admin"))
			return
		}
		panic(err)
	}

	var adminUser model.User
	err = db.WithContext(ctx).Preload("Roles").Where("email = ?", cfg.Email).First(&adminUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Admin user not found, skipping role assignment", zap.String("email", cfg.Email))
			return
		}
		panic(err)
	}

	for _, role := range adminUser.Roles {
		if role.Name == "admin" {
			logger.Warn("Admin role already assigned to user", zap.String("email", adminUser.Email))
			return
		}
	}

	err = db.WithContext(ctx).Model(&adminUser).Association("Roles").Append(adminRole)
	if err != nil {
		panic(err)
	}
	logger.Info("admin role assigned to user", zap.String("email", adminUser.Email), zap.String("role", adminRole.Name))
}
