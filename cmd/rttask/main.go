package main

import (
	"rttask/internal/config"
	"rttask/internal/domain/model"
	"rttask/internal/infrastructure/persistence/postgres"

	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoadConfig()
	logger, _ := zap.NewProduction() // вынести в MustNewLogger
	defer logger.Sync()
	db := postgres.MustNewConn(cfg.Database, logger)
	db.AutoMigrate(
		&model.Role{},
		&model.Company{},
		&model.User{},
		&model.File{},
		&model.TaskStatus{},
		&model.Task{},
		&model.Comment{},
	)

	logger.Info("config loaded", zap.String("ENV", cfg.Env))

	//router := gin.Default()
	//router.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})
	//router.Run(":8080")
}
