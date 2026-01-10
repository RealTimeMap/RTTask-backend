package dto

import (
	"mime/multipart"
	"rttask/internal/domain/model"
)

type CompanyRequest struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Avatar      *multipart.FileHeader `form:"avatar" binding:"required"`
}

type CompanyResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewCompanyResponse(company *model.Company) CompanyResponse {
	return CompanyResponse{
		ID:          company.ID,
		Name:        company.Name,
		Description: company.Description,
	}
}

func NewMultiplyCompanyResponse(companies []*model.Company) []CompanyResponse {
	response := make([]CompanyResponse, 0, len(companies))
	for _, company := range companies {
		response = append(response, NewCompanyResponse(company))
	}
	return response
}
