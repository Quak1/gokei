package service

import (
	"github.com/Quak1/gokei/internal/database/queries"
)

type Service struct {
	Hello *HelloService
}

func New(db *queries.Queries) *Service {
	return &Service{
		Hello: NewHelloService(db),
	}
}
