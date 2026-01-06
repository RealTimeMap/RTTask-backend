package valueobject

type PaginationParams struct {
	Offset int
	Limit  int
	Page   int
}

func NewPaginationParams(page, pageSize int) PaginationParams {
	offset := (page - 1) * pageSize
	return PaginationParams{
		Offset: offset,
		Limit:  pageSize,
		Page:   page,
	}
}
