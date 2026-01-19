package handlers

import (
	"net/http"
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/domain/service/invite"
	"rttask/internal/domain/valueobject"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/dto"
	"rttask/internal/transport/http/middleware"
	"rttask/internal/transport/http/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InviteHandler struct {
	service *invite.Service
	mapper  *response.ErrorMapper
	logger  *zap.Logger
}

func InitInviteHandler(g *gin.RouterGroup, service *invite.Service, logger *zap.Logger, manager security.JWTManager, mapper *response.ErrorMapper) {
	h := &InviteHandler{
		service: service,
		mapper:  mapper,
		logger:  logger,
	}
	r := g.Group("/invite")
	{
		r.POST("/", middleware.AuthMiddleware(manager, logger, mapper), h.CreateInvite)
		r.GET("/", middleware.AuthMiddleware(manager, logger, mapper), h.GetAll)
		r.PATCH("/:id", middleware.AuthMiddleware(manager, logger, mapper), h.UpdateInvite)
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
	rawData := invite.NewInviteInput(req.Description, req.RolesIDs)

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

// GetAll godoc
// @Summary Get all invite links
// @Description Get paginated list of all invite links created by the authenticated user
// @Tags invite
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} dto.InvitePaginationResponse "Successfully retrieved invite links"
// @Failure 400 {object} response.ProblemDetail "Invalid query parameters"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /invite [get]
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

// UpdateInvite godoc
// @Summary Update invite link
// @Description Update an existing invite link. Allows modifying description and associated roles.
// @Tags invite
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Invite ID"
// @Param request body dto.InviteUpdateRequest true "Invite update data"
// @Success 200 {object} dto.InviteResponse "Successfully updated invite link"
// @Failure 400 {object} response.ProblemDetail "Invalid request body or invite ID"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 404 {object} response.ProblemDetail "Invite not found"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /invite/{id} [patch]
func (h *InviteHandler) UpdateInvite(c *gin.Context) {
	inviteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		err = domainerrors.NewValidationError("invalid invite id")
		problem := h.mapper.MapError(c, err).
			WithTraceID(response.GetTraceID(c)).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	userID := response.GetUserID(c)
	traceID := response.GetTraceID(c)

	var req dto.InviteUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		problem := h.mapper.MapError(c, err).
			WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}

	updatedInvite, err := h.service.UpdateInvite(c.Request.Context(), req.ToInput(), uint(inviteID), userID)
	if err != nil {
		problem := h.mapper.MapError(c, err).
			WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	c.JSON(http.StatusOK, dto.NewInviteResponse(updatedInvite))
}
