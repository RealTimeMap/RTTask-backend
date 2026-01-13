package task

import (
	"rttask/internal/domain/valueobject"
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

type TaskFilter struct {
	valueobject.PaginationParams
	CompanyID     uint
	SortBy        string
	SortOrder     string
	ShowCompleted bool
}

func NewTaskFilter(page, pageSize int, companyID uint, sortBy, SortOrder string, showCompleted bool) TaskFilter {
	return TaskFilter{
		PaginationParams: valueobject.NewPaginationParams(page, pageSize),
		CompanyID:        companyID,
		SortBy:           sortBy,
		SortOrder:        SortOrder,
		ShowCompleted:    showCompleted,
	}
}
