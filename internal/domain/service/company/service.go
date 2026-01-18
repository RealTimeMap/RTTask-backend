package company

import (
	"context"
	"errors"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model"
	"rttask/internal/domain/model/rbac"
	"rttask/internal/domain/repository"
	"rttask/internal/domain/service/file"
	"rttask/internal/domain/valueobject"

	"go.uber.org/zap"
)

type CompanyService struct {
	companyRepo repository.CompanyRepository
	userRepo    repository.UserRepository
	fileService *file.FileService
	logger      *zap.Logger
}

func NewCompanyService(companyRepo repository.CompanyRepository, userRepo repository.UserRepository, fileService *file.FileService, logger *zap.Logger) *CompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
		userRepo:    userRepo,
		fileService: fileService,
		logger:      logger,
	}
}

func (s *CompanyService) CreateCompany(ctx context.Context, input Input, fileInput file.FileInput, userID uint) (*model.Company, error) {
	s.logger.Info("start CompanyService.CreateCompany")

	if err := s.validateUser(ctx, userID, rbac.CompanyCreate); err != nil {
		return nil, err
	}

	if err := s.validateCompanyUnique(ctx, input.Name); err != nil {
		return nil, err
	}
	logo, err := s.fileService.UploadFile(ctx, fileInput, file.CompanyProfile)
	if err != nil {
		return nil, err
	}

	company := &model.Company{
		Name:        input.Name,
		Description: input.Description,
		Avatar:      logo,
	}

	newCompany, err := s.companyRepo.Create(ctx, company)
	if err != nil {
		return nil, err
	}

	return newCompany, nil
}

func (s *CompanyService) GetAll(ctx context.Context, params valueobject.PaginationParams) ([]*model.Company, int64, error) {
	s.logger.Info("start CompanyService.GetAll")

	companies, count, err := s.companyRepo.GetAll(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	return companies, count, nil
}

func (s *CompanyService) validateCompanyUnique(ctx context.Context, name string) error {
	existCompany, err := s.companyRepo.GetByName(ctx, name)
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if errors.As(err, &notFoundErr) || notFoundErr.Type != domainerrors.ErrorTypeNotFound {
			return nil
		}
		return err
	}
	if existCompany != nil {
		return domainerrors.NewValidationError("company name is already in use")
	}
	return nil
}

func (s *CompanyService) validateUser(ctx context.Context, userID uint, permission rbac.Permission) error {
	user, err := s.userRepo.GetUserByIDWithRoles(ctx, userID)
	if err != nil {
		return err
	}
	if !user.Can(permission) {
		return domainerrors.NewForbiddenError("dont have permission")
	}
	return nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, input UpdateInput, fileInput *file.FileInput, companyID, userID uint) (*model.Company, error) {
	if err := s.validateUser(ctx, userID, rbac.CompanyUpdate); err != nil {
		return nil, err
	}

	company, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		if err := s.validateCompanyUnique(ctx, *input.Name); err != nil {
			return nil, err
		}
	}

	input.ApplyTo(company)

	if fileInput != nil {
		if company.Avatar != nil {
			if err := s.fileService.DeleteFile(ctx, company.Avatar.Path); err != nil {
				s.logger.Warn("failed to delete file", zap.Error(err), zap.String("Path", company.Avatar.Path))
			}
		}
		newFile, err := s.fileService.UploadFile(ctx, *fileInput, file.CompanyProfile)
		if err != nil {
			s.logger.Warn("failed to upload file", zap.Error(err))
			return nil, err
		}
		company.Avatar = newFile
	}

	updatedCompany, err := s.companyRepo.Update(ctx, company, companyID)
	if err != nil {
		return nil, err
	}
	return updatedCompany, nil
}
