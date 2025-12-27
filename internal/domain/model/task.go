package model

import (
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

type Task struct {
	gorm.Model
	CreatorID   uint
	Creator     User `gorm:"foreignkey:CreatorID"`
	ExecutorID  uint
	Executor    User `gorm:"foreignkey:ExecutorID"`
	Status      Status
	Priority    uint `gorm:"default:1"`
	CompanyID   uint
	Company     Company `gorm:"foreignkey:CompanyID"`
	StartAt     time.Time
	DeadlineAt  time.Time
	CompletedAt time.Time

	Files []File `gorm:"polymorphic:Entity;polymorphicValue:task"`
}
