package handlers

import (
	"net/http"
	"rttask/internal/domain/service/auth"
	"rttask/internal/domain/valueobject"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/dto"
	"rttask/internal/transport/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	manager     security.JWTManager
	authService *auth.AuthService
	logger      *zap.Logger
	mapper      *response.ErrorMapper
}

type credentials struct {
	Email    valueobject.Email
	Password valueobject.Password
}

func InitAuthHandler(g *gin.RouterGroup, manager security.JWTManager, authService *auth.AuthService) {
	authHandler := &AuthHandler{manager: manager, authService: authService}
	r := g.Group("/auth")
	{
		r.POST("/login", authHandler.Login)
		r.POST("/register", authHandler.Register)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	traceID := response.GetTraceID(c)

	if err := c.ShouldBind(&req); err != nil {
		problem := response.NewProblemDetail(
			http.StatusBadRequest,
			"Bad Request",
			"Invalid request body: "+err.Error(),
		).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}

	cred, err := h.validateCredentials(c, req.Email, req.Password)
	if err != nil {
		problem := h.mapper.MapError(c, err)
		problem.Send(c)
		return
	}

	input := auth.LoginInput{
		Email:    cred.Email,
		Password: cred.Password,
	}
	tokens, err := h.authService.Login(c.Request.Context(), input)
	if err != nil {
		problem := h.mapper.MapError(c, err)
		problem.Send(c)
		return
	}
	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cred, err := h.validateCredentials(c, req.Email, req.Password)
	if err != nil {
		problem := h.mapper.MapError(c, err)
		problem.Send(c)
		return
	}

	input := auth.RegisterInput{
		Email:      cred.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Password:   cred.Password,
		InviteLink: req.InviteLink,
	}
	user, err := h.authService.Register(c.Request.Context(), input)
	if err != nil {

		problem := h.mapper.MapError(c, err).WithTraceID(response.GetTraceID(c)).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	c.JSON(http.StatusCreated, dto.NewUserResponse(user))
}

func (h *AuthHandler) validateCredentials(c *gin.Context, email, password string) (*credentials, error) {
	traceID := response.GetTraceID(c)

	e, err := valueobject.NewEmail(email)
	if err != nil {
		h.logger.Warn("invalid email in login request",
			zap.String("traceID", traceID),
			zap.String("email", email),
		)
		return nil, err
	}
	p, err := valueobject.NewPassword(password)
	if err != nil {
		h.logger.Warn("invalid email in login request",
			zap.String("traceID", traceID),
			zap.String("email", email),
		)
		return nil, err
	}
	return &credentials{Email: e, Password: p}, nil
}
