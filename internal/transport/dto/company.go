package dto

import (
	"mime/multipart"
	"rttask/internal/domain/model"
	"rttask/internal/domain/service/company"
)

type CompanyRequest struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Avatar      *multipart.FileHeader `form:"logo" binding:"required"`
}

type CompanyResponse struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Logo        *FileResponse `json:"logo,omitempty"`
}

func NewCompanyResponse(company *model.Company) CompanyResponse {
	return CompanyResponse{
		ID:          company.ID,
		Name:        company.Name,
		Description: company.Description,
		Logo:        NewFileResponse(company.Avatar),
	}
}

func NewMultiplyCompanyResponse(companies []*model.Company) []CompanyResponse {
	response := make([]CompanyResponse, 0, len(companies))
	for _, c := range companies {
		response = append(response, NewCompanyResponse(c))
	}
	return response
}

type CompanyUpdateRequest struct {
	Name        *string               `form:"name"`
	Description *string               `form:"description"`
	Logo        *multipart.FileHeader `form:"logo"`
}

func (r *CompanyUpdateRequest) ToInput() company.UpdateInput {
	return company.UpdateInput{
		Name:        r.Name,
		Description: r.Description,
	}
}
