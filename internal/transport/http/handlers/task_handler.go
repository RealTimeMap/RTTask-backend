package handlers

import (
	"net/http"
	"rttask/internal/domain/service/file"
	"rttask/internal/domain/service/task"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/dto"
	"rttask/internal/transport/http/middleware"
	"rttask/internal/transport/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TaskHandler struct {
	service *task.TaskService
	mapper  *response.ErrorMapper
	logger  *zap.Logger
}

func InitTaskHandler(g *gin.RouterGroup, service *task.TaskService, logger *zap.Logger, manager security.JWTManager, mapper *response.ErrorMapper) {
	h := &TaskHandler{
		service: service,
		mapper:  mapper,
		logger:  logger,
	}
	r := g.Group("/task")
	{
		r.POST("/", middleware.AuthMiddleware(manager, logger, mapper), h.CreateTask)
	}
}

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
				input.File.Close()
			}
		}()
	}

	rawData := task.TaskInput{
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
	c.JSON(http.StatusCreated, newTask)
}
