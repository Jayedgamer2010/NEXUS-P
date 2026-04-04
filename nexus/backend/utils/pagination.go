package utils

import "gorm.io/gorm"

type PaginationMeta struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
}

type PaginatedResult struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

func Paginate(query *gorm.DB, page int, perPage int) *gorm.DB {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	return query.Offset(offset).Limit(perPage)
}

func GetPaginationMeta(query *gorm.DB, model interface{}, page int, perPage int) PaginationMeta {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	var total int64
	query.Session(&gorm.Session{}).Model(model).Count(&total)

	lastPage := int(total) / perPage
	if int(total)%perPage != 0 {
		lastPage++
	}

	return PaginationMeta{
		Total:       total,
		PerPage:     perPage,
		CurrentPage: page,
		LastPage:    lastPage,
	}
}

func SanitizePagination(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return page, perPage
}
