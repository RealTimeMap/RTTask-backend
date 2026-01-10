package repository

import (
	"context"
	"rttask/internal/domain/model"
	"rttask/internal/domain/valueobject"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) (*model.Task, error)
	GetUserTasks(ctx context.Context, params valueobject.PaginationParams, userID uint) ([]*model.Task, error)
	// GetAllForCompany()
	// GetByID()
}
