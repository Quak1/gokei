package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Quak1/gokei/internal/database"
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

func validateTransaction(v *validator.Validator, transaction *store.Transaction) {
	v.Check(validator.NonZero(transaction.AccountID), "account_id", "Must be provided")

	v.Check(validator.NonZero(transaction.AmountCents), "amount_cents", "Must be provided")

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

func (s *TransactionService) Create(transactionParams *store.CreateTransactionParams) (*store.Transaction, error) {
	transaction := &store.Transaction{
		AccountID:   transactionParams.AccountID,
		AmountCents: transactionParams.AmountCents,
		CategoryID:  transactionParams.CategoryID,
		Title:       transactionParams.Title,
		Attachment:  transactionParams.Attachment,
		Note:        transactionParams.Note,
	}

	v := validator.New()
	if validateTransaction(v, transaction); !v.Valid() {
		return nil, v.GetErrors()
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)
	ctx := context.Background()

	newTransaction, err := qtx.CreateTransaction(ctx, *transactionParams)
	if err != nil {
		return nil, database.HandleForeignKeyError(err)
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

	if len(data) == 0 {
		return nil, database.ErrRecordNotFound
	}

	transactions := make([]*store.Transaction, len(data))
	for i, v := range data {
		transactions[i] = &v
	}

	return transactions, nil
}

func (s *TransactionService) GetByID(id int32) (*store.Transaction, error) {
	if id < 1 {
		return nil, database.ErrRecordNotFound
	}

	transaction, err := s.queries.GetTransactionByID(context.Background(), id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &transaction, nil
}

func (s *TransactionService) DeleteByID(id int32) error {
	if id < 1 {
		return database.ErrRecordNotFound
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)
	ctx := context.Background()

	transaction, err := qtx.DeleteTransactionByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return database.ErrRecordNotFound
		default:
			return err
		}
	}

	_, err = qtx.UpdateBalance(ctx, store.UpdateBalanceParams{
		ID:           transaction.AccountID,
		BalanceCents: -transaction.AmountCents,
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

type UpdateTransactionParams struct {
	AmountCents *int64     `json:"amount_cents"`
	AccountID   *int32     `json:"account_id"`
	CategoryID  *int32     `json:"category_id"`
	Title       *string    `json:"title"`
	Date        *time.Time `json:"date"`
	Attachment  *string    `json:"attachment"`
	Note        *string    `json:"note"`
}

func (s *TransactionService) UpdateByID(id int32, updateParams *UpdateTransactionParams) (*store.Transaction, error) {
	if id < 1 {
		return nil, database.ErrRecordNotFound
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)
	ctx := context.Background()

	transaction, err := qtx.GetTransactionByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	oldAmount := transaction.AmountCents
	oldAccountID := transaction.AccountID

	if updateParams.AmountCents != nil {
		transaction.AmountCents = *updateParams.AmountCents
	}
	if updateParams.AccountID != nil {
		transaction.AccountID = *updateParams.AccountID
	}
	if updateParams.CategoryID != nil {
		transaction.CategoryID = *updateParams.CategoryID
	}
	if updateParams.Title != nil {
		transaction.Title = *updateParams.Title
	}
	if updateParams.Date != nil {
		transaction.Date = *updateParams.Date
	}
	if updateParams.Attachment != nil {
		transaction.Attachment = *updateParams.Attachment
	}
	if updateParams.Note != nil {
		transaction.Note = *updateParams.Note
	}

	v := validator.New()
	if validateTransaction(v, &transaction); !v.Valid() {
		return nil, v.GetErrors()
	}

	result, err := s.queries.UpdateTransactionById(ctx, store.UpdateTransactionByIdParams{
		ID:          transaction.ID,
		Version:     transaction.Version,
		AmountCents: transaction.AmountCents,
		AccountID:   transaction.AccountID,
		CategoryID:  transaction.CategoryID,
		Title:       transaction.Title,
		Date:        transaction.Date,
		Attachment:  transaction.Attachment,
		Note:        transaction.Note,
	})
	if err != nil {
		return nil, database.HandleForeignKeyError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, database.ErrEditConflict
	}

	_, err = qtx.UpdateBalance(ctx, store.UpdateBalanceParams{
		ID:           oldAccountID,
		BalanceCents: -oldAmount,
	})
	if err != nil {
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

	return &transaction, nil
}
