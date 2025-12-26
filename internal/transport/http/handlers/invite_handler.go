package handlers

import (
	"net/http"
	"rttask/internal/domain/service/invite"
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
	}
}

func (h *InviteHandler) CreateInvite(c *gin.Context) {
	userID := response.GetUserID(c)
	traceID := response.GetTraceID(c)

	newInvite, err := h.service.CreateInvite(c.Request.Context(), userID)

	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	c.JSON(http.StatusCreated, dto.NewInviteResponse(newInvite))

}
