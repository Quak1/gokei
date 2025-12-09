package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/pkg/validator"
)

type RecurringTransactionService struct {
	queries store.QuerierTx
	DB      *sql.DB
}

func NewRecurringTransactionService(queries store.QuerierTx, db *sql.DB) *RecurringTransactionService {
	return &RecurringTransactionService{
		queries: queries,
		DB:      db,
	}
}

func validateRecurringTransaction(v *validator.Validator, transaction *store.RecurringTransaction) {
	// TODO validate transaction
}

func (s *RecurringTransactionService) GetAll(userID int32) ([]*store.RecurringTransaction, error) {
	data, err := s.queries.GetUserRecurringTransactions(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	transactions := make([]*store.RecurringTransaction, len(data))
	for i, v := range data {
		transactions[i] = &v
	}

	return transactions, nil
}

func (s *RecurringTransactionService) Create(userID int32, params *store.CreateRecurringTransactionParams) (*store.RecurringTransaction, error) {
	if params.CategoryID == database.InitialCategoryID() {
		return nil, ErrTransactionWithInitialCategory
	}
	if params.CategoryID < 1 {
		return nil, database.ErrInvalidCategory
	}

	transaction := &store.RecurringTransaction{
		AccountID: params.AccountID,

		Title:       params.Title,
		CategoryID:  params.CategoryID,
		AmountCents: params.AmountCents,
		Note:        params.Note,

		Frequency: params.Frequency,
		Interval:  params.Interval,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,

		DayMonth: params.DayMonth,
		DayWeek:  params.DayWeek,

		MaxOccurrences: params.MaxOccurrences,
		IsActive:       params.IsActive,
	}

	v := validator.New()
	if validateRecurringTransaction(v, transaction); !v.Valid() {
		return nil, v.GetErrors()
	}

	ctx := context.Background()

	_, err := s.queries.GetAccountByID(ctx, store.GetAccountByIDParams{
		ID:     params.AccountID,
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

	newTransaction, err := s.queries.CreateRecurringTransaction(ctx, *params)
	if err != nil {
		return nil, database.HandleForeignKeyError(err)
	}

	return &newTransaction, nil
}
