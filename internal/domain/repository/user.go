package repository

import (
	"context"
	"rttask/internal/domain/model"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByIDWithRoles(ctx context.Context, id uint) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	IsUserInCompany(ctx context.Context, userID uint, companyID uint) (bool, error)
}
