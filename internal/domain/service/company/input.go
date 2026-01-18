package company

import "rttask/internal/domain/model"

type Input struct {
	Name        string
	Description string
}

type UpdateInput struct {
	Name        *string
	Description *string
}

func (i UpdateInput) ApplyTo(company *model.Company) {
	if i.Name != nil {
		company.Name = *i.Name
	}
	if i.Description != nil {
		company.Description = *i.Description
	}
}
