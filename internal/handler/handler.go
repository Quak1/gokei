package handler

import (
	"log/slog"

	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/pkg/response"
)

type Handler struct {
	Hello    *HelloHandler
	Category *CategoryHandler
}

func New(svc *service.Service, logger *slog.Logger) *Handler {
	response.SetLogger(logger)

	return &Handler{
		Hello:    NewHelloHandler(svc.Hello),
		Category: NewCategoryHandler(svc.Category),
	}
}
