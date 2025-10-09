package handler

import "github.com/Quak1/gokei/internal/service"

type Handler struct {
	Hello *HelloHandler
}

func New(svc *service.Service) *Handler {
	return &Handler{
		Hello: NewHelloHandler(svc.Hello),
	}
}
