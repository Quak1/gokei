package main

import (
	"net/http"

	"github.com/Quak1/gokei/internal/database/queries"
	"github.com/Quak1/gokei/internal/handler"
	"github.com/Quak1/gokei/internal/service"
)

func (app *application) routes(queries *queries.Queries) http.Handler {
	svc := service.New(queries)
	h := handler.New(svc)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", h.Hello.PingHandler)

	return mux
}
