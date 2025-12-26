package main

import (
	"rttask/internal/app"
	"rttask/internal/config"
	"rttask/internal/domain/model"
	"rttask/internal/domain/model/rbac"
	"rttask/internal/infrastructure/persistence/postgres"
	"rttask/internal/transport/http/handlers"
	"rttask/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoadConfig()
	logger, _ := zap.NewProduction() // вынести в MustNewLogger
	defer logger.Sync()

	db := postgres.MustNewConn(cfg.Database, logger)
	db.AutoMigrate(
		&rbac.Role{},
		&model.Company{},
		&model.User{},
		&model.File{},
		&model.TaskStatus{},
		&model.Task{},
		&model.Comment{},
		&model.InviteLink{},
	)
	logger.Info("config loaded", zap.String("ENV", cfg.Env))
	container := app.NewContainer(cfg, db, logger)
	router := gin.New()
	router.Use(middleware.TraceMiddleware())
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(gin.Logger())

	handlers.InitPermissionHandler(router.Group("/"), logger)
	handlers.InitAuthHandler(router.Group("/"), container.JWTManager, container.AuthService)
	handlers.InitInviteHandler(router.Group("/"), container.InviteService, logger, container.JWTManager, container.Mapper)

	router.Run(":8080")
}
