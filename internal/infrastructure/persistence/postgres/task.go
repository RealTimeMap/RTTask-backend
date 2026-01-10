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
	return task, nil
}

func (r *PgTaskRepository) GetUserTasks(ctx context.Context, params valueobject.PaginationParams, userID uint) ([]*model.Task, error) {
	r.logger.Info("start TaskRepository.GetUserTasks")
	var tasks []*model.Task
	err := r.db.WithContext(ctx).
		Model(&model.Task{}).
		Order("created_at desc").
		Offset(params.Offset).
		Limit(params.Limit).
		Where("user_id = ?", userID).
		Find(&tasks).Error
	if err != nil {
		return nil, MapGormError(err, "task")
	}
	return tasks, nil
}
