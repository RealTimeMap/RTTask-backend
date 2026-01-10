package handlers

import (
	"net/http"
	"rttask/internal/domain/service/auth"
	"rttask/internal/domain/service/file"
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

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "User email"
// @Param password formData string true "User password"
// @Success 200 {object} auth.Tokens "Successfully authenticated"
// @Failure 400 {object} response.ProblemDetail "Invalid request body"
// @Failure 401 {object} response.ProblemDetail "Invalid credentials"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /auth/login [post]
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

// Register godoc
// @Summary User registration
// @Description Register a new user with email, password, and invite link
// @Tags auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "User email"
// @Param password formData string true "User password"
// @Param firstName formData string true "User first name"
// @Param lastName formData string true "User last name"
// @Param inviteLink formData string true "Invite link token"
// @Success 201 {object} dto.UserResponse "Successfully registered"
// @Failure 400 {object} response.ProblemDetail "Invalid request body"
// @Failure 404 {object} response.ProblemDetail "Invite link not found"
// @Failure 409 {object} response.ProblemDetail "User already exists"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	traceID := response.GetTraceID(c)

	if err := c.ShouldBind(&req); err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	cred, err := h.validateCredentials(c, req.Email, req.Password)
	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	var fileInput *file.FileInput
	if req.Avatar != nil {
		input, err := file.NewFileInput(req.Avatar, "user", 0)
		if err != nil {
			h.logger.Error("error creating file input", zap.Error(err))
			problem := h.mapper.MapError(c, err)
			problem.Send(c)
			return
		}
		defer input.File.Close()
		fileInput = &input
	}

	input := auth.RegisterInput{
		Email:      cred.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Password:   cred.Password,
		InviteLink: req.InviteLink,
	}
	user, err := h.authService.Register(c.Request.Context(), input, fileInput)
	if err != nil {

		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
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
