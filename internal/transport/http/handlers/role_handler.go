package handlers

import (
	"net/http"
	"rttask/internal/domain/service/role"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/dto"
	"rttask/internal/transport/http/middleware"
	"rttask/internal/transport/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RoleHandler struct {
	service *role.RoleService
	logger  *zap.Logger
	mapper  *response.ErrorMapper
}

func InitRoleHandler(r *gin.RouterGroup, service *role.RoleService, logger *zap.Logger, manager security.JWTManager, mapper *response.ErrorMapper) {
	h := &RoleHandler{service: service, logger: logger, mapper: mapper}
	g := r.Group("/role")
	{
		g.POST("/", middleware.AuthMiddleware(manager, logger, mapper), h.CreateRole)
		g.GET("/permissions", h.GetAllPermissions)
	}
}

// CreateRole godoc
// @Summary Create new role
// @Description Create a new role for users
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} dto.RoleResponse "Successfully created role"
// @Failure 400 {object} response.ProblemDetail "Not valid data"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 403 {object} response.ProblemDetail "Don't have permissions"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /role [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.RoleRequest
	traceID := response.GetTraceID(c)
	userID := response.GetUserID(c)

	if err := c.ShouldBind(&req); err != nil {
		h.logger.Error("bind json error", zap.Error(err))
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}

	rawInput := role.RoleInput{Name: req.Name, Permissions: req.Permissions, UserID: userID}

	newRole, err := h.service.CreateRole(c.Request.Context(), rawInput)
	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}

	c.JSON(http.StatusCreated, dto.NewRoleResponse(newRole))
}

// GetAllPermissions godoc
// @Summary Systems permissions
// @Description Get all system permissions for role.
// @Tags roles
// @Produce json
// @Success 200 {object} []dto.PermissionGroupDTO "Successfully registered"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /role/permissions [get]
func (h *RoleHandler) GetAllPermissions(c *gin.Context) {
	permResponse := dto.NewGroupedPermissions()
	c.JSON(http.StatusOK, permResponse)
}
