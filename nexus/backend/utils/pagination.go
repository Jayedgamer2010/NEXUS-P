package utils

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Meta struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	From        int   `json:"from"`
	To          int   `json:"to"`
}

type Paginated struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

func GetPage(c *fiber.Ctx) int {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	return page
}

func GetPerPage(c *fiber.Ctx) int {
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return perPage
}

func Paginate(query *gorm.DB, page, perPage int) *gorm.DB {
	offset := (page - 1) * perPage
	return query.Offset(offset).Limit(perPage)
}

func BuildMeta(total int64, page, perPage int) Meta {
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	from := (page-1)*perPage + 1
	to := page * perPage
	if int64(to) > total {
		to = int(total)
	}
	return Meta{
		Total:       total,
		PerPage:     perPage,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        from,
		To:          to,
	}
}
