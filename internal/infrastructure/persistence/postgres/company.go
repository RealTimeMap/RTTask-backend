package postgres

import (
	"context"
	"rttask/internal/domain/model"
	"rttask/internal/domain/repository"
	"rttask/internal/domain/valueobject"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PgCompanyRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewPgCompanyRepository(db *gorm.DB, logger *zap.Logger) repository.CompanyRepository {
	return &PgCompanyRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PgCompanyRepository) Create(ctx context.Context, company *model.Company) (*model.Company, error) {
	r.logger.Info("start CompanyRepository.Create")
	err := r.db.WithContext(ctx).Create(company).Error
	if err != nil {
		return nil, MapGormError(err, "company")
	}
	return company, nil
}

func (r *PgCompanyRepository) GetByName(ctx context.Context, name string) (*model.Company, error) {
	var company *model.Company
	err := r.db.WithContext(ctx).First(&company, "name = ?", name).Error
	if err != nil {
		return nil, MapGormError(err, "company")
	}
	return company, nil
}
func (r *PgCompanyRepository) GetAll(ctx context.Context, params valueobject.PaginationParams) ([]*model.Company, int64, error) {
	r.logger.Info("start CompanyRepository.GetAll")
	var companies []*model.Company
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Company{}).Offset(params.Offset).Limit(params.Limit).Find(&companies).Count(&count).Error
	if err != nil {
		return nil, 0, MapGormError(err, "company")
	}
	return companies, count, nil
}

func (r *PgCompanyRepository) GetByID(ctx context.Context, id uint) (*model.Company, error) {
	r.logger.Info("start CompanyRepository.GetByID")
	var company model.Company
	err := r.db.WithContext(ctx).First(&company, "id = ?", id).Error
	if err != nil {
		return nil, MapGormError(err, "company")
	}
	return &company, nil
}
