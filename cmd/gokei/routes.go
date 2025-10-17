package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/ping", app.handler.Hello.PingHandler)

	mux.HandleFunc("GET /v1/categories", app.handler.Category.GetAll)
	mux.HandleFunc("POST /v1/categories", app.handler.Category.Create)
	mux.HandleFunc("GET /v1/categories/{categoryID}", app.handler.Category.GetByID)
	mux.HandleFunc("DELETE /v1/categories/{categoryID}", app.handler.Category.DeleteByID)
	mux.HandleFunc("PUT /v1/categories/{categoryID}", app.handler.Category.UpdateByID)

	mux.HandleFunc("GET /v1/accounts", app.handler.Account.GetAll)
	mux.HandleFunc("POST /v1/accounts", app.handler.Account.Create)
	mux.HandleFunc("GET /v1/accounts/{accountID}/transactions", app.handler.Transaction.GetAccountTransactions)

	mux.HandleFunc("GET /v1/transactions", app.handler.Transaction.GetAll)
	mux.HandleFunc("POST /v1/transactions", app.handler.Transaction.Create)

	return mux
}
