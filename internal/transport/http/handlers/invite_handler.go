package handlers

import (
	"net/http"
	"rttask/internal/domain/service/invite"
	"rttask/internal/domain/valueobject"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/dto"
	"rttask/internal/transport/http/middleware"
	"rttask/internal/transport/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InviteHandler struct {
	service *invite.InviteService
	mapper  *response.ErrorMapper
	logger  *zap.Logger
}

func InitInviteHandler(g *gin.RouterGroup, service *invite.InviteService, logger *zap.Logger, manager security.JWTManager, mapper *response.ErrorMapper) {
	h := &InviteHandler{
		service: service,
		mapper:  mapper,
		logger:  logger,
	}
	r := g.Group("/invite")
	{
		r.POST("/", middleware.AuthMiddleware(manager, logger, mapper), h.CreateInvite)
		r.GET("/", middleware.AuthMiddleware(manager, logger, mapper), h.GetAll)
	}
}

// CreateInvite godoc
// @Summary Create invite link
// @Description Create a new invite link for user registration
// @Tags invite
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} dto.InviteResponse "Successfully created invite link"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /invite [post]
func (h *InviteHandler) CreateInvite(c *gin.Context) {
	var req dto.InviteRequest
	userID := response.GetUserID(c)
	traceID := response.GetTraceID(c)

	if err := c.ShouldBind(&req); err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	h.logger.Info("ids", zap.Any("ids", req.RolesIDs))
	rawData := invite.NewInviteInput(req.RolesIDs)

	newInvite, err := h.service.CreateInvite(c.Request.Context(), rawData, userID)

	if err != nil {
		problem := h.mapper.MapError(c, err).
			WithTraceID(traceID).
			WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	c.JSON(http.StatusCreated, dto.NewInviteResponse(newInvite))

}

func (h *InviteHandler) GetAll(c *gin.Context) {
	var params dto.PaginationRequest
	params.Default()

	traceID := response.GetTraceID(c)
	userID := response.GetUserID(c)

	if err := c.ShouldBindQuery(&params); err != nil {
		problem := h.mapper.MapError(c, err).
			WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}

	validParams := valueobject.PaginationParams{Page: params.Page, Offset: params.Offset(), Limit: params.Limit()}

	invites, err := h.service.GetAllInvites(c.Request.Context(), userID, validParams)
	if err != nil {
		problem := h.mapper.MapError(c, err).
			WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	c.JSON(http.StatusOK, dto.NewPaginationResponse(invites, params, 100))
}
