package model

import (
	domainerrors "rttask/internal/domain/errors"
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	CreatedStatus    Status = "Новый"
	InWorkStatus     Status = "В работе"
	InProgressStatus Status = "В доработке"
	CompletedStatus  Status = "Выполнена"
	ImmediateStatus  Status = "Срочная"
)

func ParseStatus(s string) (Status, error) {
	switch Status(s) {
	case CreatedStatus, InWorkStatus, InProgressStatus, CompletedStatus, ImmediateStatus:
		return Status(s), nil
	}
	return "", domainerrors.NewValidationError("invalid status")
}

type Task struct {
	gorm.Model
	CreatorID   uint
	Creator     User `gorm:"foreignkey:CreatorID"`
	ExecutorID  uint
	Executor    User `gorm:"foreignkey:ExecutorID"`
	Title       string
	Description string
	Status      Status
	Priority    uint `gorm:"default:1"`
	CompanyID   uint
	Company     Company `gorm:"foreignkey:CompanyID"`
	StartAt     time.Time
	DeadlineAt  time.Time
	CompletedAt time.Time

	Files []*File `gorm:"type:jsonb;serializer:json"`
}

type TaskStatusGroup struct {
	Status Status
	Tasks  []*Task
	Count  int64
}

type GroupedTasksByStatus struct {
	Groups     []*TaskStatusGroup
	TotalCount int64
}

type TaskNotifier interface {
	NotifyNewTask(task *Task)
}
