package service

import (
	"github.com/Quak1/gokei/internal/database/queries"
)

type Service struct {
	Hello    *HelloService
	Category *CategoryService
}

func New(db *queries.Queries) *Service {
	return &Service{
		Hello:    NewHelloService(db),
		Category: NewCategoryService(db),
	}
}
