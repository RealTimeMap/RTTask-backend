package repository

import (
	"context"
	"rttask/internal/domain/model"
	"rttask/internal/domain/valueobject"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *model.Company) (*model.Company, error)
	GetByName(ctx context.Context, name string) (*model.Company, error)
	GetAll(ctx context.Context, params valueobject.PaginationParams) ([]*model.Company, int64, error)
	GetByID(ctx context.Context, id uint) (*model.Company, error)
}
