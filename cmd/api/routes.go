package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/ping", app.handler.Hello.PingHandler)
	mux.HandleFunc("POST /v1/categories", app.handler.Category.Create)

	return mux
}
