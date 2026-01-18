package task

import (
	"rttask/internal/domain/model"
	"time"
)

type CreateInput struct {
	StartAt     time.Time
	DeadlineAt  time.Time
	Title       string
	Description string
	Priority    uint
	ExecutorID  uint
	CompanyID   uint
}

type UpdateInput struct {
	Title          *string
	Description    *string
	Priority       *uint
	ExecutorID     *uint
	CompanyID      *uint
	DeleteFilesIDs []string
	StartAt        *time.Time
	DeadlineAt     *time.Time
}

func (i UpdateInput) ApplyTo(task *model.Task) {
	if i.Title != nil {
		task.Title = *i.Title
	}
	if i.Description != nil {
		task.Description = *i.Description
	}
	if i.Priority != nil {
		task.Priority = *i.Priority
	}
	if i.ExecutorID != nil {
		task.ExecutorID = *i.ExecutorID
	}
	if i.StartAt != nil {
		task.StartAt = *i.StartAt
	}
	if i.DeadlineAt != nil {
		task.DeadlineAt = *i.DeadlineAt
	}
}
