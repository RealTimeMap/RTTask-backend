package main

import (
	"context"
	"rttask/internal/app"
	"rttask/internal/config"
	"rttask/internal/domain/model"
	"rttask/internal/domain/model/rbac"
	"rttask/internal/infrastructure/persistence/postgres"
	"rttask/internal/scripts"
	"rttask/internal/transport/http/handlers"
	"rttask/internal/transport/http/middleware"
	"time"

	_ "rttask/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title RTTask API
// @version 1.0
// @description Task management system API documentation
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /
// @schemes https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
		&model.Task{},
		&model.Comment{},
		&model.InviteLink{},
	)
	logger.Info("config loaded", zap.String("ENV", cfg.Env))

	container := app.NewContainer(cfg, db, logger)

	scripts.CreateAdminIfNotExists(context.Background(), cfg.Admin, logger, container.UserRepository, container.Hasher)

	router := gin.Default()
	router.Use(middleware.TraceMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://rt-task-frontend.vercel.app", "https://realtimemap.ru", "http://localhost:5173", "http://localhost:1420", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Trace-Id"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.JSON(404, gin.H{"error": "Not found"})
	})

	router.NoMethod(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.JSON(405, gin.H{"error": "Method not allowed"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	handlers.InitAuthHandler(router.Group("/"), container.JWTManager, container.AuthService)
	handlers.InitInviteHandler(router.Group("/"), container.InviteService, logger, container.JWTManager, container.Mapper)
	handlers.InitRoleHandler(router.Group("/"), container.RoleService, logger, container.JWTManager, container.Mapper)

	router.Run(":8081")
}
