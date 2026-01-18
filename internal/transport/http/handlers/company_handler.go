package handlers

import (
	"net/http"
	"rttask/internal/domain/errors"
	"rttask/internal/domain/service/company"
	"rttask/internal/domain/service/file"
	"rttask/internal/domain/valueobject"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/dto"
	"rttask/internal/transport/http/middleware"
	"rttask/internal/transport/http/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CompanyHandler struct {
	service *company.CompanyService
	mapper  *response.ErrorMapper
	logger  *zap.Logger
}

func InitCompanyHandler(g *gin.RouterGroup, service *company.CompanyService, logger *zap.Logger, manager security.JWTManager, mapper *response.ErrorMapper) {
	h := &CompanyHandler{
		service: service,
		mapper:  mapper,
		logger:  logger,
	}
	r := g.Group("/company")
	{
		r.POST("/", middleware.AuthMiddleware(manager, logger, mapper), h.CreateCompany)
		r.GET("/", middleware.AuthMiddleware(manager, logger, mapper), h.GetCompanies)
		r.PATCH("/:id", middleware.AuthMiddleware(manager, logger, mapper), h.UpdateCompany)
	}
}

// CreateCompany godoc
// @Summary Create new company
// @Description Create a new company with avatar image
// @Tags companies
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Company name"
// @Param description formData string true "Company description"
// @Param avatar formData file true "Company avatar image"
// @Success 201 {object} dto.CompanyResponse "Successfully created company"
// @Failure 400 {object} response.ProblemDetail "Invalid request body"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /company [post]
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req dto.CompanyRequest
	userID := response.GetUserID(c)
	traceID := response.GetTraceID(c)

	if err := c.ShouldBind(&req); err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}

	fileInput, err := file.NewFileInput(req.Avatar, "company", userID)
	if err != nil {
		h.logger.Error("failed to create file input", zap.Error(err))
		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	defer fileInput.File.Close()

	companyInput := company.Input{
		Name:        req.Name,
		Description: req.Description,
	}

	// Вызвать сервис
	newCompany, err := h.service.CreateCompany(c.Request.Context(), companyInput, fileInput, userID)
	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID).WithInstance(c.Request.URL.Path)
		problem.Send(c)
		return
	}
	res := dto.NewCompanyResponse(newCompany)
	c.JSON(http.StatusCreated, res)
}

// GetCompanies godoc
// @Summary Get companies list
// @Description Get paginated list of all companies
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} dto.CompanyPaginationResponse "Successfully retrieved companies"
// @Failure 400 {object} response.ProblemDetail "Invalid query parameters"
// @Failure 401 {object} response.ProblemDetail "Unauthorized - invalid or missing token"
// @Failure 500 {object} response.ProblemDetail "Internal server error"
// @Router /company [get]
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	var params dto.PaginationRequest
	params.Default()

	traceID := response.GetTraceID(c)

	if err := c.ShouldBind(&params); err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}
	validParams := valueobject.NewPaginationParams(params.Page, params.PageSize)

	companies, count, err := h.service.GetAll(c.Request.Context(), validParams)
	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}

	companiesResponse := dto.NewMultiplyCompanyResponse(companies)

	c.JSON(http.StatusOK, dto.NewPaginationResponse(companiesResponse, params, count))
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	userID := response.GetUserID(c)
	traceID := response.GetTraceID(c)
	companyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		err = errors.NewValidationError("invalid company ID")
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}

	var req dto.CompanyUpdateRequest

	if err := c.ShouldBind(&req); err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}

	var fileInput *file.FileInput
	if req.Logo != nil {
		input, err := file.NewFileInput(req.Logo, "company", userID)
		if err != nil {
			h.logger.Error("failed to create file input", zap.Error(err))
			problem := h.mapper.MapError(c, err).WithTraceID(traceID)
			problem.Send(c)
			return
		}
		fileInput = &input
		defer fileInput.File.Close()
	}

	updatedCompany, err := h.service.UpdateCompany(c.Request.Context(), req.ToInput(), fileInput, uint(companyID), userID)
	if err != nil {
		problem := h.mapper.MapError(c, err).WithTraceID(traceID)
		problem.Send(c)
		return
	}
	c.JSON(http.StatusOK, dto.NewCompanyResponse(updatedCompany))
}
