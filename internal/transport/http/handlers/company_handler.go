package handlers

import (
	"net/http"
	"rttask/internal/domain/service/company"
	"rttask/internal/domain/service/file"
	"rttask/internal/domain/valueobject"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/dto"
	"rttask/internal/transport/http/middleware"
	"rttask/internal/transport/http/response"

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
	}
}

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

	companyInput := company.CompanyInput{
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

	c.JSON(http.StatusCreated, newCompany)
}

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
