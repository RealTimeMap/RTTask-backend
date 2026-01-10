package task

import (
	"time"
)

type TaskInput struct {
	StartAt     time.Time
	DeadlineAt  time.Time
	Title       string
	Description string
	Priority    uint
	ExecutorID  uint
	CompanyID   uint
}
