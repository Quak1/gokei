package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/ping", app.handler.Hello.PingHandler)

	mux.HandleFunc("POST /v1/users", app.handler.User.Create)
	mux.HandleFunc("GET /v1/users/{userID}", app.handler.User.GetByID)
	mux.HandleFunc("PUT /v1/users/{userID}", app.handler.User.UpdateByID)
	mux.HandleFunc("POST /v1/users/login", app.handler.User.Login)
	mux.Handle("DELETE /v1/users/{userID}", app.handler.User.Authenticate(http.HandlerFunc(app.handler.User.DeleteByID)))

	mux.HandleFunc("GET /v1/categories", app.handler.Category.GetAll)
	mux.HandleFunc("POST /v1/categories", app.handler.Category.Create)
	mux.HandleFunc("GET /v1/categories/{categoryID}", app.handler.Category.GetByID)
	mux.HandleFunc("PUT /v1/categories/{categoryID}", app.handler.Category.UpdateByID)
	mux.HandleFunc("DELETE /v1/categories/{categoryID}", app.handler.Category.DeleteByID)

	mux.HandleFunc("GET /v1/accounts", app.handler.Account.GetAll)
	mux.HandleFunc("POST /v1/accounts", app.handler.Account.Create)
	mux.HandleFunc("GET /v1/accounts/{accountID}", app.handler.Account.GetByID)
	mux.HandleFunc("GET /v1/accounts/{accountID}/balance", app.handler.Account.GetSumBalance)
	mux.HandleFunc("GET /v1/accounts/{accountID}/transactions", app.handler.Transaction.GetAccountTransactions)
	mux.HandleFunc("PUT /v1/accounts/{accountID}", app.handler.Account.UpdateByID)
	mux.HandleFunc("DELETE /v1/accounts/{accountID}", app.handler.Account.DeleteByID)

	mux.HandleFunc("GET /v1/transactions", app.handler.Transaction.GetAll)
	mux.HandleFunc("POST /v1/transactions", app.handler.Transaction.Create)
	mux.HandleFunc("GET /v1/transactions/{transactionID}", app.handler.Transaction.GetByID)
	mux.HandleFunc("PUT /v1/transactions/{transactionID}", app.handler.Transaction.UpdateByID)
	mux.HandleFunc("DELETE /v1/transactions/{transactionID}", app.handler.Transaction.DeleteByID)

	return mux
}
