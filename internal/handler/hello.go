package handler

import (
	"net/http"

	"github.com/Quak1/gokei/internal/service"
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

	w.Write([]byte(msg))
}
