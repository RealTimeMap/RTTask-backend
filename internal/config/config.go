package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Database struct {
	Host         string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port         int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User         string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password     string `yaml:"password" env:"DB_PASSWORD"`
	DBName       string `yaml:"dbName" env:"DB_NAME" env-required:"true"`
	SSLMode      string `yaml:"sslMode" env:"DB_SSL_MODE" env-default:"disable"`
	MaxOpenConn  int    `yaml:"maxOpenConn" env:"DB_MAX_OPEN_CONN" env-default:"25"`
	ConnLifeTime int    `yaml:"connLifeTime" env:"DB_CONN_LIFE_TIME" env-default:"5"`
	MaxIdleConn  int    `yaml:"MaxIdleConn" env:"DB_MAX_IDLE_CONN" env-default:"10"`
}

func (db Database) ConnMaxLifetime() time.Duration {
	return time.Duration(db.ConnLifeTime) * time.Minute
}

func (db Database) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode,
	)
}

type JWT struct {
	Secret               string `yaml:"secretKey" env:"JWT_SECRET_KEY" env-required:"true"`
	AccessTokenDuration  int    `yaml:"accessTokenDuration" env:"JWT_ACCESS_TOKEN_DURATION" env-default:"24"`
	RefreshTokenDuration int    `yaml:"refreshTokenDuration" env:"JWT_REFRESH_TOKEN_DURATION" env-default:"1000"`
}

func (J JWT) AccessTokenTimeDuration() time.Duration {
	return time.Duration(J.AccessTokenDuration) * time.Minute
}

func (J JWT) RefreshTokenTimeDuration() time.Duration {
	return time.Duration(J.RefreshTokenDuration) * time.Minute
}

type Admin struct {
	Email    string `yaml:"email" env:"ADMIN_EMAIL" env-default:"admin@admin.ru"`
	Password string `yaml:"password" env:"ADMIN_PASSWORD" env-default:"admin1admin"`
}

type Config struct {
	Env      string   `env:"ENV" env-default:"local"`
	Database Database `yaml:"database"`
	JWT      JWT      `yaml:"jwt"`
	Admin    Admin    `yaml:"admin"`
}

func MustLoadConfig() Config {
	var cfg Config
	path := "./config/config.yaml"
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
