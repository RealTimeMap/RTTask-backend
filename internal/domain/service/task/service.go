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
	"rttask/internal/domain/valueobject"
	"slices"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	taskRepo    repository.TaskRepository
	userRepo    repository.UserRepository
	companyRepo repository.CompanyRepository

	notifier model.TaskNotifier

	fileService *file.FileService
	logger      *zap.Logger
}

func NewTaskService(taskRepo repository.TaskRepository, userRepo repository.UserRepository, companyRepo repository.CompanyRepository, fileService *file.FileService, logger *zap.Logger) *Service {
	return &Service{
		taskRepo:    taskRepo,
		userRepo:    userRepo,
		companyRepo: companyRepo,
		fileService: fileService,
		logger:      logger,
	}
}

func (s *Service) SetNotifier(notifier model.TaskNotifier) {
	s.notifier = notifier
}

func (s *Service) CreateTask(ctx context.Context, input CreateInput, filesInput []file.FileInput, userID uint) (*model.Task, error) {
	// Валидация пользователя кем поставлена задача
	if err := s.validateUser(ctx, userID, rbac.TaskCreate, rbac.TaskAssign); err != nil {
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
			uploadedFile, err := s.fileService.UploadFile(ctx, fileInput, file.TaskProfile)
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

	go s.notifier.NotifyNewTask(newTask)

	return newTask, nil
}

// GetTasks Получение задач с пагинацией и фильтрами
func (s *Service) GetTasks(ctx context.Context, params valueobject.TaskFilterList, userID uint) ([]*model.Task, int64, error) {
	if err := s.validateUser(ctx, userID, rbac.TaskView, rbac.TaskList); err != nil {
		return nil, 0, err
	}
	if err := s.validateCompany(ctx, params.CompanyID); err != nil {
		return nil, 0, err
	}

	tasks, total, err := s.taskRepo.GetTasks(ctx, params)
	if err != nil {
		s.logger.Error("failed to get tasks", zap.Error(err))
		return nil, 0, err
	}
	return tasks, total, nil

}

// GetUserTasks получение пользовательских задач в группированом ввиде
func (s *Service) GetUserTasks(ctx context.Context, userID uint) (*model.GroupedTasksByStatus, error) {
	if err := s.validateUser(ctx, userID, rbac.TaskView, rbac.TaskList); err != nil {
		return nil, err
	}

	tasks, err := s.taskRepo.GetUserTasks(ctx, userID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// ChangeTaskStatus изменение статуса задачи
func (s *Service) ChangeTaskStatus(ctx context.Context, taskID uint, status string, userID uint) (*model.Task, error) {
	if err := s.validateUser(ctx, userID, rbac.TaskChangeStatus); err != nil {
		return nil, err
	}
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.ExecutorID != userID || task.CreatorID != userID {
		return nil, domainerrors.NewForbiddenError("task does not own user")
	}

	newStatus, err := model.ParseStatus(status)
	if err != nil {
		return nil, err
	}

	task.Status = newStatus

	updatedTask, err := s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil

}

// DeleteTask Удаление задачи
func (s *Service) DeleteTask(ctx context.Context, taskID, userID uint) error {
	if err := s.validateUser(ctx, userID, rbac.TaskDelete); err != nil {
		return err
	}
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}
	if err = s.taskRepo.Delete(ctx, taskID); err != nil {
		return err
	}
	return nil
}

// UpdateTask обновление задачи
func (s *Service) UpdateTask(ctx context.Context, input UpdateInput, filesInput []file.FileInput, taskID, userID uint) (*model.Task, error) {
	if err := s.validateUser(ctx, userID, rbac.TaskUpdate); err != nil {
		return nil, err
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Применяем изменения полей
	input.ApplyTo(task)

	// Удаляем файлы
	if len(input.DeleteFilesIDs) > 0 {
		task.Files = s.removeFiles(ctx, task.Files, input.DeleteFilesIDs)
	}

	// Добавляем новые файлы
	for _, fileInput := range filesInput {
		uploadedFile, err := s.fileService.UploadFile(ctx, fileInput, file.TaskProfile)
		if err != nil {
			return nil, err
		}
		task.Files = append(task.Files, uploadedFile)
	}

	return s.taskRepo.Update(ctx, task)
}

func (s *Service) removeFiles(ctx context.Context, files []*model.File, deleteIDs []string) []*model.File {
	result := make([]*model.File, 0, len(files))
	for _, f := range files {
		if slices.Contains(deleteIDs, f.ID) {
			// Удаляем файл
			if err := s.fileService.DeleteFile(ctx, f.Path); err != nil {
				s.logger.Error("failed to delete file", zap.String("path", f.Path), zap.Error(err))
			}
			continue
		}
		result = append(result, f)
	}
	return result
}

func (s *Service) validateUser(ctx context.Context, userID uint, permissions ...rbac.Permission) error {
	user, err := s.userRepo.GetUserByIDWithRoles(ctx, userID)
	if err != nil {
		return err
	}
	if !user.CanAll(permissions...) {
		return domainerrors.NewForbiddenError("dont have permission")
	}
	return nil
}

func (s *Service) validateExecutor(ctx context.Context, executorID uint, companyID uint) error {
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

func (s *Service) validateCompany(ctx context.Context, companyID uint) error {
	_, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		var notFoundErr *domainerrors.DomainError
		if errors.As(err, &notFoundErr) || notFoundErr.Type != domainerrors.ErrorTypeNotFound {
			s.logger.Info("company not found", zap.Uint("companyID", companyID))
			return err
		}
	}
	return nil
}

func (s *Service) checkUserInCompany(ctx context.Context, userID uint, companyID uint) error {
	inCompany, err := s.userRepo.IsUserInCompany(ctx, userID, companyID)
	if err != nil {
		return err
	}
	if !inCompany {
		return domainerrors.NewValidationError(fmt.Sprintf("User %d not in company", userID))
	}
	return nil

}

func (s *Service) validateTimeRange(startAt time.Time, deadlineAt time.Time) error {
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

func (s *Service) validatePriority(value uint) error {
	if value > 5 || value == 0 {
		return domainerrors.NewValidationError("priority must be between 0 and 5")
	}
	return nil
}
