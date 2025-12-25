package http

import (
	"net/http"
	"rttask/internal/transport/dto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PermissionsHandler struct {
	logger *zap.Logger
}

func InitPermissionHandler(g *gin.RouterGroup, logger *zap.Logger) {
	h := &PermissionsHandler{logger: logger}
	r := g.Group("/permissions")
	{
		r.GET("/all", h.GetAllPermissions)
	}
}

func (h *PermissionsHandler) GetAllPermissions(c *gin.Context) {
	response := dto.NewGroupedPermissions()
	c.JSON(http.StatusOK, response)
}
