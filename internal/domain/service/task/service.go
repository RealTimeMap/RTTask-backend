package task

import (
	"context"
	"errors"
	"fmt"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/model"
	"rttask/internal/domain/model/rbac"
	"rttask/internal/domain/repository"
	"rttask/internal/domain/service/file"
	"time"

	"go.uber.org/zap"
)

type TaskService struct {
	taskRepo    repository.TaskRepository
	userRepo    repository.UserRepository
	companyRepo repository.CompanyRepository
	fileService *file.FileService
	logger      *zap.Logger
}

func NewTaskService(taskRepo repository.TaskRepository, userRepo repository.UserRepository, companyRepo repository.CompanyRepository, fileService *file.FileService, logger *zap.Logger) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		userRepo:    userRepo,
		companyRepo: companyRepo,
		fileService: fileService,
		logger:      logger,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, input TaskInput, filesInput []file.FileInput, userID uint) (*model.Task, error) {
	// Валидация пользователя кем поставлена задача
	if err := s.validateCreator(ctx, userID, rbac.TaskCreate, rbac.TaskAssign); err != nil {
		s.logger.Error("failed to validate user", zap.Error(err))
		return nil, err
	}

	// Валидация компании
	if err := s.validateCompany(ctx, input.CompanyID); err != nil {
		s.logger.Error("failed to validate company", zap.Error(err))
		return nil, err
	}

	// Валидация исполнителя, что он может менять статус задачи и принадлежит к той же компании
	if err := s.validateExecutor(ctx, input.ExecutorID, input.CompanyID); err != nil {
		s.logger.Error("failed to validate executor", zap.Error(err))
		return nil, err
	}

	// Валидация временных промежутков
	if err := s.validateTimeRange(input.StartAt, input.DeadlineAt); err != nil {
		s.logger.Error("failed to validate time", zap.Error(err))
		return nil, err
	}

	task := &model.Task{
		Title:       input.Title,
		Description: input.Description,
		StartAt:     input.StartAt,
		DeadlineAt:  input.DeadlineAt,
		Priority:    input.Priority,
		ExecutorID:  input.ExecutorID,
		CompanyID:   input.CompanyID,
		CreatorID:   userID,
		Status:      model.CreatedStatus,
	}

	if len(filesInput) > 0 {
		for _, fileInput := range filesInput {
			uploadedFile, err := s.fileService.UploadFile(ctx, fileInput, file.TaskProfile) // нужно создать
			if err != nil {
				return nil, err
			}
			task.Files = append(task.Files, uploadedFile)
		}
	}

	newTask, err := s.taskRepo.Create(ctx, task)
	if err != nil {
		s.logger.Error("failed to create task", zap.Error(err))
		return nil, err
	}

	return newTask, nil
}

func (s *TaskService) validateCreator(ctx context.Context, userID uint, permissions ...rbac.Permission) error {
	user, err := s.userRepo.GetUserByIDWithRoles(ctx, userID)
	if err != nil {
		return err
	}
	if !user.CanAll(permissions...) {
		return domainerrors.NewForbiddenError("dont have permission")
	}
	return nil
}

func (s *TaskService) validateExecutor(ctx context.Context, executorID uint, companyID uint) error {
	executor, err := s.userRepo.GetUserByIDWithRoles(ctx, executorID)
	if err != nil {
		return err
	}
	if !executor.CanAll(rbac.TaskView, rbac.TaskChangeStatus) {
		return domainerrors.NewValidationError("Executor not allowed to view tasks and change his status")
	}
	if err := s.checkUserInCompany(ctx, executorID, companyID); err != nil {
		return err
	}
	return nil
}

func (s *TaskService) validateCompany(ctx context.Context, companyID uint) error {
	_, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if !errors.As(err, &notFoundErr) || notFoundErr.Type != domainerrors.ErrorTypeNotFound {
			s.logger.Error("failed to check company existence",
				zap.Uint("companyID", companyID),
				zap.Error(err),
			)
			return err
		}
	}
	return nil
}

func (s *TaskService) checkUserInCompany(ctx context.Context, userID uint, companyID uint) error {
	inCompany, err := s.userRepo.IsUserInCompany(ctx, userID, companyID)
	if err != nil {
		return err
	}
	if !inCompany {
		return domainerrors.NewValidationError(fmt.Sprintf("User %d not in company", userID))
	}
	return nil

}

func (s *TaskService) validateTimeRange(startAt time.Time, deadlineAt time.Time) error {
	now := time.Now()

	if deadlineAt.Before(startAt) {
		return domainerrors.NewValidationError("deadline is before startAt")
	}

	maxPastDays := 7
	minAllowedDate := now.AddDate(0, 0, -maxPastDays)

	if startAt.Before(minAllowedDate) {
		return domainerrors.NewValidationError("start date too far in the past")
	}

	maxFutureDays := 365
	maxAllowedDays := now.AddDate(0, 0, maxFutureDays)

	if deadlineAt.After(maxAllowedDays) {
		return domainerrors.NewValidationError("deadline is too far in the future")
	}

	return nil

}
