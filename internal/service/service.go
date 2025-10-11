package service

import (
	"github.com/Quak1/gokei/internal/database/store"
)

type Service struct {
	Hello    *HelloService
	Category *CategoryService
}

func New(queries *store.Queries) *Service {
	return &Service{
		Hello:    NewHelloService(queries),
		Category: NewCategoryService(queries),
	}
}
