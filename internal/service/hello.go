package service

import (
	"github.com/Quak1/gokei/internal/database/queries"
)

type HelloService struct {
	db *queries.Queries
}

func NewHelloService(db *queries.Queries) *HelloService {
	return &HelloService{
		db: db,
	}
}

func (s *HelloService) GetMessage() string {
	return "Hello, world!"
}
