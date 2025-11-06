package main

import (
	"net/http"

	"github.com/Quak1/gokei/internal/middleware"
	"github.com/Quak1/gokei/internal/service"
)

func (app *application) routes(s *service.Service) http.Handler {
	mw := middleware.New(s)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/ping", app.handler.Hello.PingHandler)

	mux.HandleFunc("POST /v1/users", app.handler.User.Create)
	mux.Handle("GET /v1/users/{userID}", mw.Authenticate(http.HandlerFunc(app.handler.User.GetByID)))
	mux.Handle("PUT /v1/users/{userID}", mw.Authenticate(http.HandlerFunc(app.handler.User.UpdateByID)))
	mux.Handle("DELETE /v1/users/{userID}", mw.Authenticate(http.HandlerFunc(app.handler.User.DeleteByID)))

	mux.HandleFunc("POST /v1/auth/login", app.handler.Auth.Login)

	mux.HandleFunc("GET /v1/categories", app.handler.Category.GetAll)
	mux.HandleFunc("POST /v1/categories", app.handler.Category.Create)
	mux.HandleFunc("GET /v1/categories/{categoryID}", app.handler.Category.GetByID)
	mux.HandleFunc("PUT /v1/categories/{categoryID}", app.handler.Category.UpdateByID)
	mux.HandleFunc("DELETE /v1/categories/{categoryID}", app.handler.Category.DeleteByID)

	mux.Handle("GET /v1/accounts", mw.Authenticate(http.HandlerFunc(app.handler.Account.GetAll)))
	mux.Handle("POST /v1/accounts", mw.Authenticate(http.HandlerFunc(app.handler.Account.Create)))
	mux.Handle("GET /v1/accounts/{accountID}", mw.Authenticate(http.HandlerFunc(app.handler.Account.GetByID)))
	mux.Handle("GET /v1/accounts/{accountID}/balance", mw.Authenticate(http.HandlerFunc(app.handler.Account.GetSumBalance)))
	mux.Handle("GET /v1/accounts/{accountID}/transactions", mw.Authenticate(http.HandlerFunc(app.handler.Transaction.GetAccountTransactions)))
	mux.Handle("PUT /v1/accounts/{accountID}", mw.Authenticate(http.HandlerFunc(app.handler.Account.UpdateByID)))
	mux.Handle("DELETE /v1/accounts/{accountID}", mw.Authenticate(http.HandlerFunc(app.handler.Account.DeleteByID)))

	mux.Handle("GET /v1/transactions", mw.Authenticate(http.HandlerFunc(app.handler.Transaction.GetAll)))
	mux.HandleFunc("POST /v1/transactions", app.handler.Transaction.Create)
	mux.Handle("GET /v1/transactions/{transactionID}", mw.Authenticate(http.HandlerFunc(app.handler.Transaction.GetByID)))
	mux.Handle("PUT /v1/transactions/{transactionID}", mw.Authenticate(http.HandlerFunc(app.handler.Transaction.UpdateByID)))
	mux.Handle("DELETE /v1/transactions/{transactionID}", mw.Authenticate(http.HandlerFunc(app.handler.Transaction.DeleteByID)))

	return mux
}
