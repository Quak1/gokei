package service

import (
	"github.com/Quak1/gokei/internal/database"
)

type Service struct {
	Hello       *HelloService
	Category    *CategoryService
	Account     *AccountService
	Transaction *TransactionService
}

func New(db *database.DB) *Service {
	return &Service{
		Hello:       NewHelloService(db.Queries),
		Category:    NewCategoryService(db.Queries),
		Account:     NewAccountService(db.Queries),
		Transaction: NewTransactionService(db.Queries, db.Connection),
	}
}
