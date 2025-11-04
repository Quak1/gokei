package service

import (
	"github.com/Quak1/gokei/internal/database"
)

type Service struct {
	Hello       *HelloService
	Category    *CategoryService
	Account     *AccountService
	Transaction *TransactionService
	User        *UserService
	Token       *TokenService
}

func New(db *database.DB) *Service {
	tokenService := NewTokenService(db.Queries)

	return &Service{
		Hello:       NewHelloService(db.Queries),
		Category:    NewCategoryService(db.Queries),
		Account:     NewAccountService(db.Queries, db.Connection),
		Transaction: NewTransactionService(db.Queries, db.Connection),
		User:        NewUserService(db.Queries, tokenService),
		Token:       tokenService,
	}
}
