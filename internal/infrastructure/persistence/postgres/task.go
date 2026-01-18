package postgres

import (
	"context"
	"rttask/internal/domain/model"
	"rttask/internal/domain/repository"
	"rttask/internal/domain/valueobject"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PgTaskRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewPgTaskRepository(db *gorm.DB, logger *zap.Logger) repository.TaskRepository {
	return &PgTaskRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PgTaskRepository) Create(ctx context.Context, task *model.Task) (*model.Task, error) {
	r.logger.Info("start TaskRepository.Create")
	err := r.db.WithContext(ctx).Create(&task).Error
	if err != nil {
		return nil, MapGormError(err, "task")
	}

	return r.GetByID(ctx, task.ID)
}

func (r *PgTaskRepository) GetUserTasks(ctx context.Context, userID uint) (*model.GroupedTasksByStatus, error) {
	r.logger.Info("start TaskRepository.GetUserTasks")

	type StatusCount struct {
		Status model.Status
		Count  int64
	}

	var statusCounts []StatusCount
	err := r.db.WithContext(ctx).
		Model(&model.Task{}).
		Select("status, COUNT(*) as count").
		Where("executor_id = ?", userID).
		Group("status").
		Find(&statusCounts).Error
	if err != nil {
		return nil, MapGormError(err, "task")
	}

	result := &model.GroupedTasksByStatus{
		Groups:     make([]*model.TaskStatusGroup, 0),
		TotalCount: 0,
	}

	for _, sc := range statusCounts {
		var tasks []*model.Task
		err := r.db.WithContext(ctx).
			Model(&model.Task{}).
			Preload("Creator").
			Preload("Executor").
			Preload("Company").
			Where("executor_id = ? AND status = ?", userID, sc.Status).
			Find(&tasks).Error
		if err != nil {
			return nil, MapGormError(err, "task")
		}

		result.Groups = append(result.Groups, &model.TaskStatusGroup{
			Status: sc.Status,
			Tasks:  tasks,
			Count:  sc.Count,
		})
		result.TotalCount += sc.Count
	}

	return result, nil
}

func (r *PgTaskRepository) GetTasks(ctx context.Context, params valueobject.TaskFilterList) ([]*model.Task, int64, error) {
	r.logger.Info("start TaskRepository.GetTasks")
	var tasks []*model.Task
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Task{}).
		Preload("Creator").
		Preload("Executor").
		Preload("Company").
		Order(params.GetOrderBy()).
		Offset(params.Offset).
		Limit(params.Limit).
		Where("company_id = ?", params.CompanyID).
		Count(&count).
		Find(&tasks).Error
	if err != nil {
		return nil, 0, MapGormError(err, "task")
	}
	return tasks, count, nil
}

func (r *PgTaskRepository) GetByID(ctx context.Context, id uint) (*model.Task, error) {
	r.logger.Info("start TaskRepository.GetByID")
	var task model.Task
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Executor").
		Preload("Company").
		First(&task, id).Error
	if err != nil {
		return nil, MapGormError(err, "task")
	}
	return &task, nil
}

func (r *PgTaskRepository) Update(ctx context.Context, task *model.Task) (*model.Task, error) {
	r.logger.Info("start TaskRepository.Update")
	err := r.db.WithContext(ctx).Model(&model.Task{}).Where("id = ?", task.ID).Save(task).Error
	if err != nil {
		return nil, MapGormError(err, "task")
	}
	return r.GetByID(ctx, task.ID)
}

func (r *PgTaskRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("start TaskRepository.Delete")
	err := r.db.WithContext(ctx).Model(&model.Task{}).Where("id = ?", id).Delete(&model.Task{}).Error
	if err != nil {
		return MapGormError(err, "task")
	}
	return nil
}

func (r *PgTaskRepository) ParticipantUpdate(ctx context.Context, data map[string]interface{}, id uint) (*model.Task, error) {
	r.logger.Info("start TaskRepository.ParticipantUpdate")
	err := r.db.WithContext(ctx).Model(&model.Task{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, MapGormError(err, "task")
	}
	return r.GetByID(ctx, id)
}
