package dto

import (
	"mime/multipart"
	"time"
)

type TaskParams struct {
	PaginationRequest
	CompanyID int `json:"companyId"`
}

type TaskRequest struct {
	Title       string                  `form:"title" binding:"required"`
	Description string                  `form:"description" binding:"required"`
	Priority    uint                    `form:"priority" binding:"required"`
	ExecutorID  uint                    `form:"executorId" binding:"required"`
	CompanyID   uint                    `form:"companyId" binding:"required"`
	StartAt     time.Time               `form:"startAt" binding:"required"`
	DeadlineAt  time.Time               `form:"deadlineAt" binding:"required"`
	Files       []*multipart.FileHeader `form:"files"`
}

type TaskResponse struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    uint      `json:"priority"`
	ExecutorID  uint      `json:"executorId"`
	CompanyID   uint      `json:"companyId"`
	StartAt     time.Time `json:"startAt"`
	DeadlineAt  time.Time `json:"deadlineAt"`
}
