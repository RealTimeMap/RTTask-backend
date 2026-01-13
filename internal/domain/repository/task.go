package repository

import (
	"context"
	"rttask/internal/domain/model"
	"rttask/internal/domain/valueobject"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) (*model.Task, error)
	GetUserTasks(ctx context.Context, userID uint) (*model.GroupedTasksByStatus, error)
	GetTasks(ctx context.Context, params valueobject.TaskFilterList) ([]*model.Task, int64, error)
	// GetByID()
}
