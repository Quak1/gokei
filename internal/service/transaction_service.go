package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/pkg/validator"
)

var (
	ErrDeleteInitialTransaction       = errors.New("Can't delete initial transaction")
	ErrRefundInitialTransaction       = errors.New("Can't refund initial transaction")
	ErrTransactionWithInitialCategory = errors.New("Can't create transaction with initial category")
)

type TransactionService struct {
	queries store.QuerierTx
	DB      *sql.DB
}

func NewTransactionService(queries store.QuerierTx, db *sql.DB) *TransactionService {
	return &TransactionService{
		queries: queries,
		DB:      db,
	}
}

func validateTransaction(v *validator.Validator, transaction *store.Transaction) {
	v.Check(validator.NonZero(transaction.AccountID), "account_id", "Must be provided")

	// v.Check(validator.NonZero(transaction.AmountCents), "amount_cents", "Must be provided")

	v.Check(validator.NonZero(transaction.CategoryID), "category_id", "Must be provided")

	v.Check(validator.NonZero(transaction.Title), "title", "Must be provided")

	// v.Check(validator.NonZero(transaction.Attachment), "attachment", "Must be provided")
	// TODO extend validation

	// v.Check(validator.NonZero(transaction.Note), "note", "Must be provided")
}

func (s *TransactionService) GetAll(userID int32) ([]*store.Transaction, error) {
	data, err := s.queries.GetAllTransactions(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	transactions := make([]*store.Transaction, len(data))
	for i, v := range data {
		transactions[i] = &v.Transaction
	}

	return transactions, nil
}

func (s *TransactionService) Create(userID int32, transactionParams *store.CreateTransactionParams) (*store.Transaction, error) {
	if transactionParams.CategoryID == database.InitialCategoryID() {
		return nil, ErrTransactionWithInitialCategory
	}
	if transactionParams.CategoryID < 1 {
		return nil, database.ErrRecordNotFound
	}

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

	_, err = qtx.GetAccountByID(ctx, store.GetAccountByIDParams{
		ID:     transactionParams.AccountID,
		UserID: userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	newTransaction, err := qtx.CreateTransaction(ctx, *transactionParams)
	if err != nil {
		return nil, database.HandleForeignKeyError(err)
	}

	_, err = qtx.UpdateBalance(ctx, store.UpdateBalanceParams{
		ID:           transaction.AccountID,
		BalanceCents: transaction.AmountCents,
		UserID:       userID,
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &newTransaction, nil
}

func (s *TransactionService) GetAllTRansactionsForAccountID(accountID, userID int32) ([]*store.Transaction, error) {
	data, err := s.queries.GetTransactionsByAccountID(context.Background(), store.GetTransactionsByAccountIDParams{
		AccountID: accountID,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, database.ErrRecordNotFound
	}

	transactions := make([]*store.Transaction, len(data))
	for i, v := range data {
		transactions[i] = &v.Transaction
	}

	return transactions, nil
}

func (s *TransactionService) GetByID(transactionID, userID int32) (*store.Transaction, error) {
	if transactionID < 1 || userID < 1 {
		return nil, database.ErrRecordNotFound
	}

	transaction, err := s.queries.GetTransactionByID(context.Background(), store.GetTransactionByIDParams{
		ID:     transactionID,
		UserID: userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &transaction.Transaction, nil
}

func (s *TransactionService) DeleteByID(transactionID, userID int32) error {
	if transactionID < 1 || userID < 1 {
		return database.ErrRecordNotFound
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)
	ctx := context.Background()

	t, err := qtx.GetTransactionByID(ctx, store.GetTransactionByIDParams{
		ID:     transactionID,
		UserID: userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return database.ErrRecordNotFound
		default:
			return err
		}
	}

	transaction := t.Transaction
	if t.Transaction.CategoryID == database.InitialCategoryID() {
		return ErrDeleteInitialTransaction
	}

	result, err := qtx.DeleteTransactionByID(ctx, store.DeleteTransactionByIDParams{
		ID:     transactionID,
		UserID: userID,
	})
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	_, err = qtx.UpdateBalance(ctx, store.UpdateBalanceParams{
		ID:           transaction.AccountID,
		BalanceCents: -transaction.AmountCents,
		UserID:       userID,
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

func (s *TransactionService) UpdateByID(transactionID, userID int32, updateParams *UpdateTransactionParams) (*store.Transaction, error) {
	if transactionID < 1 || userID < 1 {
		return nil, database.ErrRecordNotFound
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)
	ctx := context.Background()

	t, err := qtx.GetTransactionByID(ctx, store.GetTransactionByIDParams{
		ID:     transactionID,
		UserID: userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	transaction := t.Transaction

	oldAmount := transaction.AmountCents
	oldAccountID := transaction.AccountID

	if updateParams.AmountCents != nil {
		transaction.AmountCents = *updateParams.AmountCents
	}
	if updateParams.AccountID != nil {
		transaction.AccountID = *updateParams.AccountID
	}
	if updateParams.CategoryID != nil {
		if *updateParams.CategoryID == database.InitialCategoryID() {
			return nil, ErrTransactionWithInitialCategory
		}
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

	_, err = qtx.GetAccountByID(ctx, store.GetAccountByIDParams{
		ID:     transaction.AccountID,
		UserID: userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
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
		UserID:      userID,
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
		UserID:       userID,
	})
	if err != nil {
		return nil, err
	}

	_, err = qtx.UpdateBalance(ctx, store.UpdateBalanceParams{
		ID:           transaction.AccountID,
		BalanceCents: transaction.AmountCents,
		UserID:       userID,
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

type RefundTransactionParams struct {
	Reason *string `json:"reason"`
}

func (s *TransactionService) RefundByID(transactionID, userID int32, params *RefundTransactionParams) (*store.Transaction, error) {
	if transactionID < 1 || userID < 1 {
		return nil, database.ErrRecordNotFound
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)
	ctx := context.Background()

	t, err := qtx.GetTransactionByID(ctx, store.GetTransactionByIDParams{
		ID:     transactionID,
		UserID: userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	transaction := t.Transaction
	if t.Transaction.CategoryID == database.InitialCategoryID() {
		return nil, ErrRefundInitialTransaction
	}

	refundTransaction, err := qtx.CreateTransaction(ctx, store.CreateTransactionParams{
		AccountID:   transaction.AccountID,
		AmountCents: -transaction.AmountCents,
		CategoryID:  transaction.CategoryID,
		Title:       fmt.Sprintf("[REFUND #%d] %s", transaction.ID, transaction.Title),
		Attachment:  transaction.Attachment,
		Note:        transaction.Note,
	})
	if err != nil {
		return nil, err
	}

	_, err = qtx.UpdateBalance(ctx, store.UpdateBalanceParams{
		ID:           refundTransaction.AccountID,
		BalanceCents: refundTransaction.AmountCents,
		UserID:       userID,
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &refundTransaction, nil
}
