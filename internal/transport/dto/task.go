package dto

import (
	"mime/multipart"
	"rttask/internal/domain/model"
	"rttask/internal/domain/service/task"
	"time"
)

type TaskFilterParams struct {
	PaginationRequest
	CompanyID     uint   `form:"companyId"`
	SortBy        string `form:"sortBy"`
	SortOrder     string `form:"sortOrder"`
	ShowCompleted bool   `form:"showCompleted"`
}

// SetDefaults устанавливает значения по умолчанию для фильтров
func (p *TaskFilterParams) SetDefaults() {
	// Пагинация
	if p.Page == 0 || p.PageSize == 0 {
		p.Default()
	}

	// Сортировка по умолчанию: createdAt desc (от новых к старым)
	if p.SortBy == "" {
		p.SortBy = "createdAt"
	}
	if p.SortOrder == "" {
		p.SortOrder = "desc"
	}

	if p.SortBy != "priority" && p.SortBy != "createdAt" {
		p.SortBy = "createdAt"
	}
	if p.SortOrder != "asc" && p.SortOrder != "desc" {
		p.SortOrder = "desc"
	}
}

// GetOrderBy возвращает строку для ORDER BY в SQL запросе
func (p *TaskFilterParams) GetOrderBy() string {
	order := "DESC"
	if p.SortOrder == "asc" {
		order = "ASC"
	}

	if p.SortBy == "priority" {
		return "priority " + order
	}
	return "created_at " + order
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
	ID          uint            `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Priority    uint            `json:"priority"`
	Creator     UserResponse    `json:"creator"`
	Executor    UserResponse    `json:"executor"`
	Company     CompanyResponse `json:"company"`
	Files       []FileResponse  `json:"files,omitempty"`
	StartAt     time.Time       `json:"startAt"`
	DeadlineAt  time.Time       `json:"deadlineAt"`
	CompletedAt *time.Time      `json:"completedAt,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func NewTaskResponse(task *model.Task) TaskResponse {
	var completedAt *time.Time
	if !task.CompletedAt.IsZero() {
		completedAt = &task.CompletedAt
	}

	return TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		Priority:    task.Priority,
		Creator:     NewUserResponse(&task.Creator),
		Executor:    NewUserResponse(&task.Executor),
		Company:     NewCompanyResponse(&task.Company),
		Files:       NewMultiFileResponse(task.Files),
		StartAt:     task.StartAt,
		DeadlineAt:  task.DeadlineAt,
		CompletedAt: completedAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func NewMultiTaskResponse(tasks []*model.Task) []TaskResponse {
	response := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, NewTaskResponse(task))
	}
	return response
}

// TaskStatusGroupResponse - группа задач с одним статусом
type TaskStatusGroupResponse struct {
	Status string         `json:"status"`
	Tasks  []TaskResponse `json:"tasks"`
	Count  int64          `json:"count"`
}

// GroupedTasksByStatusResponse - результат группировки задач по статусам
type GroupedTasksByStatusResponse struct {
	Groups     []TaskStatusGroupResponse `json:"groups"`
	TotalCount int64                     `json:"totalCount"`
}

// NewTaskStatusGroupResponse - конвертер из model в DTO для группы
func NewTaskStatusGroupResponse(group *model.TaskStatusGroup) TaskStatusGroupResponse {
	return TaskStatusGroupResponse{
		Status: string(group.Status),
		Tasks:  NewMultiTaskResponse(group.Tasks),
		Count:  group.Count,
	}
}

// NewGroupedTasksByStatusResponse - конвертер из model в DTO для группировки
func NewGroupedTasksByStatusResponse(grouped *model.GroupedTasksByStatus) GroupedTasksByStatusResponse {
	groups := make([]TaskStatusGroupResponse, 0, len(grouped.Groups))
	for _, group := range grouped.Groups {
		groups = append(groups, NewTaskStatusGroupResponse(group))
	}

	return GroupedTasksByStatusResponse{
		Groups:     groups,
		TotalCount: grouped.TotalCount,
	}
}

type TaskStatusUpdateRequest struct {
	Status string `form:"status" json:"status" binding:"required"`
}

type TaskUpdateRequest struct {
	Title          *string                 `form:"title"`
	Description    *string                 `form:"description"`
	Priority       *uint                   `form:"priority"`
	ExecutorID     *uint                   `form:"executorId"`
	CompanyID      *uint                   `form:"companyId"`
	StartAt        *time.Time              `form:"startAt"`
	DeadlineAt     *time.Time              `form:"deadlineAt"`
	DeleteFilesIDs []string                `form:"deleteFilesIds"`
	Files          []*multipart.FileHeader `form:"files"`
}

func (r *TaskUpdateRequest) ToUpdateInput() task.UpdateInput {
	return task.UpdateInput{
		Title:          r.Title,
		Description:    r.Description,
		Priority:       r.Priority,
		ExecutorID:     r.ExecutorID,
		CompanyID:      r.CompanyID,
		StartAt:        r.StartAt,
		DeadlineAt:     r.DeadlineAt,
		DeleteFilesIDs: r.DeleteFilesIDs,
	}
}
