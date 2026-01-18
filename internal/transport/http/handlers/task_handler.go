package handlers

import (
	"net/http"
	"rttask/internal/domain/errors"
	"rttask/internal/domain/service/file"
	"rttask/internal/domain/service/task"
	"rttask/internal/domain/valueobject"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/dto"
	"rttask/internal/transport/http/middleware"
	"rttask/internal/transport/http/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TaskHandler struct {
	service *task.Service
	mapper  *response.ErrorMapper
	logger  *zap.Logger
}

func InitTaskHandler(g *gin.RouterGroup, service *task.Service, logger *zap.Logger, manager security.JWTManager, mapper *response.ErrorMapper) {
	h := &TaskHandler{
		service: service,
		mapper:  mapper,
		logger:  logger,
	}
	r := g.Group("/task")
	{
		r.POST("/", middleware.AuthMiddleware(manager, logger, mapper), h.CreateTask)
		r.GET("/", middleware.AuthMiddleware(manager, logger, mapper), h.GetTasks)
		r.GET("/my", middleware.AuthMiddleware(manager, logger, mapper), h.GetUserTasks)
		r.PATCH("/:id/status", middleware.AuthMiddleware(manager, logger, mapper), h.UpdateTaskStatus)
		r.DELETE("/:id", middleware.AuthMiddleware(manager, logger, mapper), h.DeleteTask)
		r.PATCH("/:id", middleware.AuthMiddleware(manager, logger, mapper), h.UpdateTask)
	}
}

// CreateTask godoc
// @Summary Create new task
// @Description Create a new task with files attachments
// @Tags tasks
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Task title"
// @Param description formData string true "Task description"
// @Param priority formData int true "Task priority"
// @Param executorId formData int true "Executor user ID"
// @Param companyId formData int true "Company ID"
// @Param startAt formData string true "Start date" format(date-time)
// @Param deadlineAt formData string true "Deadline date" format(date-time)
// @Param files formData file false "Task files"
// @Success 201 {object} dto.TaskResponse "Successfully created task"
// @Failure 400 {object} response.ProblemDetail "Invalid request body"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /task [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.TaskRequest
	userID := response.GetUserID(c)
	traceID := response.GetTraceID(c)

	if err := c.ShouldBind(&req); err != nil {
		h.logger.Error("failed to bind request", zap.Error(err))
		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}

	var fileInputs []file.FileInput
	if len(req.Files) > 0 {
		fileInputs = make([]file.FileInput, 0, len(req.Files))
		for _, fileHeader := range req.Files {
			input, err := file.NewFileInput(fileHeader, "task", userID)
			if err != nil {
				h.logger.Error("failed to create file input", zap.Error(err))
				problem := h.mapper.MapError(c, err).WithTraceID(traceID)
				problem.Send(c)
				return
			}
			fileInputs = append(fileInputs, input)
		}
		defer func() {
			for _, input := range fileInputs {
				_ = input.File.Close()
			}
		}()
	}

	rawData := task.CreateInput{
		Title:       req.Title,
		Description: req.Description,
		CompanyID:   req.CompanyID,
		ExecutorID:  req.ExecutorID,
		StartAt:     req.StartAt,
		DeadlineAt:  req.DeadlineAt,
		Priority:    req.Priority,
	}

	newTask, err := h.service.CreateTask(c.Request.Context(), rawData, fileInputs, userID)
	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}
	c.JSON(http.StatusCreated, dto.NewTaskResponse(newTask))
}

// GetTasks godoc
// @Summary Get tasks list
// @Description Get paginated and filtered list of tasks for a company
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param companyId query int true "Company ID"
// @Param sortBy query string false "Sort field (createdAt, priority)" default(createdAt)
// @Param sortOrder query string false "Sort order (asc, desc)" default(desc)
// @Param showCompleted query bool false "Show completed tasks" default(false)
// @Success 200 {object} dto.TaskPaginationResponse "Successfully retrieved tasks"
// @Failure 400 {object} response.ProblemDetail "Invalid query parameters"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /task [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	var params dto.TaskFilterParams

	userID := response.GetUserID(c)
	traceID := response.GetTraceID(c)

	if err := c.ShouldBindQuery(&params); err != nil {
		h.logger.Error("failed to bind request", zap.Error(err))
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}

	params.SetDefaults()
	inputParams := valueobject.NewTaskFilterList(params.Page, params.PageSize, params.SortBy, params.SortOrder, params.CompanyID, params.ShowCompleted)

	tasks, total, err := h.service.GetTasks(c.Request.Context(), inputParams, userID)
	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}

	res := dto.NewMultiTaskResponse(tasks)
	c.JSON(http.StatusOK, dto.NewPaginationResponse(res, params.PaginationRequest, total))
}

// GetUserTasks godoc
// @Summary Get user tasks grouped by status
// @Description Get all tasks assigned to the authenticated user, grouped by status
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.GroupedTasksByStatusResponse "Successfully retrieved grouped tasks"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /task/my [get]
func (h *TaskHandler) GetUserTasks(c *gin.Context) {

	traceID := response.GetTraceID(c)
	userID := response.GetUserID(c)

	tasks, err := h.service.GetUserTasks(c.Request.Context(), userID)
	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}
	res := dto.NewGroupedTasksByStatusResponse(tasks)
	c.JSON(http.StatusOK, res)
}

// UpdateTaskStatus godoc
// @Summary Update task status
// @Description Change the status of an existing task. Only the task executor or creator can change the status. Requires 'task:changeStatus' permission. Valid statuses: "Новый", "В работе", "В доработке", "Выполнена", "Срочная"
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param request body dto.TaskStatusUpdateRequest true "New status"
// @Success 200 {object} dto.TaskResponse "Successfully updated task status"
// @Failure 400 {object} response.ProblemDetail "Invalid request body or invalid status transition"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 403 {object} response.ProblemDetail "Forbidden - not executor/creator or missing permission"
// @Failure 404 {object} response.ProblemDetail "Task not found"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /task/{id}/status [patch]
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	traceID := response.GetTraceID(c)
	userID := response.GetUserID(c)

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}

	var req dto.TaskStatusUpdateRequest

	if err := c.ShouldBind(&req); err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}

	updatedTask, err := h.service.ChangeTaskStatus(c.Request.Context(), uint(taskID), req.Status, userID)
	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	c.JSON(http.StatusOK, dto.NewTaskResponse(updatedTask))

}

// DeleteTask godoc
// @Summary Delete task
// @Description Delete an existing task by ID. Requires 'task:delete' permission. Only the task creator or users with delete permission can delete tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 204 "Successfully deleted task"
// @Failure 400 {object} response.ProblemDetail "Invalid task ID"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 403 {object} response.ProblemDetail "Forbidden - missing permission"
// @Failure 404 {object} response.ProblemDetail "Task not found"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /task/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	traceID := response.GetTraceID(c)
	userID := response.GetUserID(c)
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		problem := h.mapper.MapError(c, errors.NewValidationError("task id should be an integer")).WithTraceID(traceID)
		problem.Send(c)
		return
	}
	if err := h.service.DeleteTask(c.Request.Context(), uint(taskID), userID); err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateTask godoc
// @Summary Update task
// @Description Update an existing task. Supports partial updates - only provided fields will be changed. Files can be added via 'files' field and removed via 'deleteFilesIds' field.
// @Tags tasks
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param title formData string false "Task title"
// @Param description formData string false "Task description"
// @Param priority formData int false "Task priority (1-5)"
// @Param executorId formData int false "Executor user ID"
// @Param startAt formData string false "Start date" format(date-time)
// @Param deadlineAt formData string false "Deadline date" format(date-time)
// @Param files formData file false "New files to attach"
// @Param deleteFilesIds formData []string false "IDs of files to delete"
// @Success 200 {object} dto.TaskResponse "Successfully updated task"
// @Failure 400 {object} response.ProblemDetail "Invalid request body"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 403 {object} response.ProblemDetail "Forbidden - missing permission"
// @Failure 404 {object} response.ProblemDetail "Task not found"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /task/{id} [patch]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	traceID := response.GetTraceID(c)
	userID := response.GetUserID(c)
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		problem := h.mapper.MapError(c, errors.NewValidationError("task id should be an integer")).WithTraceID(traceID)
		problem.Send(c)
		return
	}

	var req dto.TaskUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}

	var fileInputs []file.FileInput
	if len(req.Files) > 0 {
		fileInputs = make([]file.FileInput, 0, len(req.Files))
		for _, fileHeader := range req.Files {
			input, err := file.NewFileInput(fileHeader, "task", userID)
			if err != nil {
				h.logger.Error("failed to create file input", zap.Error(err))
				problem := h.mapper.MapError(c, err).WithTraceID(traceID)
				problem.Send(c)
				return
			}
			fileInputs = append(fileInputs, input)
		}
		defer func() {
			for _, input := range fileInputs {
				_ = input.File.Close()
			}
		}()
	}

	updatedTask, err := h.service.UpdateTask(c.Request.Context(), req.ToUpdateInput(), fileInputs, uint(taskID), userID)
	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}
	c.JSON(http.StatusOK, dto.NewTaskResponse(updatedTask))

}
