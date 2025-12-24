package postgres

import (
	"rttask/internal/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustNewConn(cfg config.Database, logger *zap.Logger) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		panic("failed to connect database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		panic("failed to connect database")
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime())

	return db
}
