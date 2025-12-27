package handlers

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

// GetAllPermissions godoc
// @Summary Systems permissions
// @Description Get all system permissions for roles.
// @Tags permissions
// @Produce json
// @Success 200 {object} []dto.PermissionGroupDTO "Successfully registered"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /permissions/all [get]
func (h *PermissionsHandler) GetAllPermissions(c *gin.Context) {
	permResponse := dto.NewGroupedPermissions()
	c.JSON(http.StatusOK, permResponse)
}
