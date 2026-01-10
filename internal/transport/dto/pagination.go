package dto

import "math"

type PaginationRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationRequest) Limit() int {
	return p.PageSize
}

func (p *PaginationRequest) Default() {
	p.Page = 1
	p.PageSize = 10
}

type PaginationResponse[T any] struct {
	Items      []T   `json:"items"`
	TotalPages int   `json:"totalPages"`
	Total      int64 `json:"total"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

func NewPaginationResponse[T any](items []T, params PaginationRequest, total int64) PaginationResponse[T] {
	totalPages := 0
	if params.PageSize > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(params.PageSize)))
	}

	return PaginationResponse[T]{
		Items:      items,
		TotalPages: totalPages,
		Total:      total,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}
}
