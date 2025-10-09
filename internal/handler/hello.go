package handler

import (
	"net/http"

	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/pkg/response"
)

type HelloHandler struct {
	service *service.HelloService
}

func NewHelloHandler(svc *service.HelloService) *HelloHandler {
	return &HelloHandler{
		service: svc,
	}
}

func (h *HelloHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	msg := h.service.GetMessage()

	err := response.OK(w, response.Envelope{"message": msg})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}
