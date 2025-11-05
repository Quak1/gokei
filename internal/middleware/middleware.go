package middleware

import (
	"github.com/Quak1/gokei/internal/service"
)

type Middleware struct {
	service *service.Service
}

func New(service *service.Service) *Middleware {
	return &Middleware{
		service: service,
	}
}
