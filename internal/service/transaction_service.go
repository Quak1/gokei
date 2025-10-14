package service

import (
	"context"
	"database/sql"

	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/pkg/validator"
)

type TransactionService struct {
	queries *store.Queries
	DB      *sql.DB
}

func NewTransactionService(queries *store.Queries, db *sql.DB) *TransactionService {
	return &TransactionService{
		queries: queries,
		DB:      db,
	}
}

func (s *TransactionService) validateTransaction(v *validator.Validator, transaction *store.CreateTransactionParams) {
	v.Check(validator.NonZero(transaction.AccountID), "account_id", "Must be provided")

	v.Check(validator.NonZero(transaction.AmountCents), "amount", "Must be provided")

	v.Check(validator.NonZero(transaction.CategoryID), "category_id", "Must be provided")

	v.Check(validator.NonZero(transaction.Title), "title", "Must be provided")

	// v.Check(validator.NonZero(transaction.Attachment), "attachment", "Must be provided")
	// TODO extend validation

	// v.Check(validator.NonZero(transaction.Note), "note", "Must be provided")
}

func (s *TransactionService) GetAll() ([]*store.Transaction, error) {
	data, err := s.queries.GetAllTransactions(context.Background())
	if err != nil {
		return nil, err
	}

	transactions := make([]*store.Transaction, len(data))
	for i, v := range data {
		transactions[i] = &v
	}

	return transactions, nil
}

func (s *TransactionService) Create(transaction *store.CreateTransactionParams) (*store.Transaction, error) {
	v := validator.New()
	if s.validateTransaction(v, transaction); !v.Valid() {
		return nil, v.GetErrors()
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)
	ctx := context.Background()

	newTransaction, err := qtx.CreateTransaction(ctx, *transaction)
	if err != nil {
		// TODO handle account or category doesnt exist
		// pq: insert or update on table \"transactions\" violates foreign key constraint \"transactions_category_id_fkey\"
		// pq: insert or update on table \"transactions\" violates foreign key constraint \"transactions_account_id_fkey\"
		return nil, err
	}

	_, err = qtx.UpdateBalance(ctx, store.UpdateBalanceParams{
		ID:           transaction.AccountID,
		BalanceCents: transaction.AmountCents,
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// TODO return category name?
	return &newTransaction, nil
}

func (s *TransactionService) GetAllTRansactionsForAccountID(accountID int) ([]*store.Transaction, error) {
	data, err := s.queries.GetTransactionsByAccountID(context.Background(), int32(accountID))
	if err != nil {
		return nil, err
	}

	transactions := make([]*store.Transaction, len(data))
	for i, v := range data {
		transactions[i] = &v
	}

	return transactions, nil
}
