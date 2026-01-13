package valueobject

type TaskFilterBase struct {
	PaginationParams
	SortBy    string
	SortOrder string
}

func (f TaskFilterBase) GetOrderBy() string {
	direction := "DESC"
	if f.SortOrder == "asc" {
		direction = "ASC"
	}
	field := "created_at"
	if f.SortBy == "priority" {
		field = "priority"
	}
	return field + " " + direction
}

func NewTaskFilterBase(page, pageSize int, sortBy, sortOrder string) TaskFilterBase {
	if sortBy != "priority" && sortBy != "createdAt" {
		sortBy = "createdAt"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	return TaskFilterBase{
		PaginationParams: NewPaginationParams(page, pageSize),
		SortBy:           sortBy,
		SortOrder:        sortOrder,
	}
}

type TaskFilterUser struct {
	TaskFilterBase
	UserID uint
}

func NewTaskFilterUser(page, pageSize int, sortBy, sortOrder string, userID uint) TaskFilterUser {
	return TaskFilterUser{
		TaskFilterBase: NewTaskFilterBase(page, pageSize, sortBy, sortOrder),
		UserID:         userID,
	}
}

type TaskFilterList struct {
	TaskFilterBase
	CompanyID     uint
	ShowCompleted bool
}

func NewTaskFilterList(page, pageSize int, sortBy, sortOrder string, companyID uint, showCompleted bool) TaskFilterList {
	return TaskFilterList{
		TaskFilterBase: NewTaskFilterBase(page, pageSize, sortBy, sortOrder),
		CompanyID:      companyID,
		ShowCompleted:  showCompleted,
	}
}
