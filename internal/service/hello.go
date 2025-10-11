package service

import (
	"github.com/Quak1/gokei/internal/database/store"
)

type HelloService struct {
	queries *store.Queries
}

func NewHelloService(queries *store.Queries) *HelloService {
	return &HelloService{
		queries: queries,
	}
}

func (s *HelloService) GetMessage() string {
	return "Hello, world!"
}
