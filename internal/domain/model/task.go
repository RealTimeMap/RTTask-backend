package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	CreatorID   uint
	Creator     User `gorm:"foreignkey:CreatorID"`
	ExecutorID  uint
	Executor    User `gorm:"foreignkey:ExecutorID"`
	StatusID    uint
	Status      TaskStatus `gorm:"foreignkey:StatusID"`
	CompanyID   uint
	Company     Company `gorm:"foreignkey:CompanyID"`
	StartAt     time.Time
	DeadlineAt  time.Time
	CompletedAt time.Time

	Files []File `gorm:"polymorphic:Entity;polymorphicValue:task"`
}
