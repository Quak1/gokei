package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/pkg/validator"
)

type AccountService struct {
	queries store.QuerierTx
	DB      *sql.DB
}

func NewAccountService(queries store.QuerierTx, db *sql.DB) *AccountService {
	return &AccountService{
		queries: queries,
		DB:      db,
	}
}

func validateAccount(v *validator.Validator, account *store.Account) {
	v.Check(validator.NonZero(account.Name), "name", "Must be provided")
	v.Check(validator.MaxLength(account.Name, 50), "name", "Must not be more than 50 bytes long")

	v.Check(validator.NonZero(account.Type), "type", "Must be provided")
	v.Check(validator.PermittedValue(account.Type, "debit", "cash", "credit"), "type", "Invalid account type. Valid types are credit, debit, and cash")
}

func (s *AccountService) GetAll(userID int32) ([]*store.Account, error) {
	data, err := s.queries.GetUserAccounts(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	accounts := make([]*store.Account, len(data))
	for i, a := range data {
		accounts[i] = &a
	}

	return accounts, nil
}

func (s *AccountService) Create(accountParams *store.CreateAccountParams) (*store.Account, error) {
	account := &store.Account{
		Name: accountParams.Name,
		Type: accountParams.Type,
	}

	v := validator.New()
	if validateAccount(v, account); !v.Valid() {
		return nil, v.GetErrors()
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)
	ctx := context.Background()

	newAccount, err := qtx.CreateAccount(ctx, *accountParams)
	if err != nil {
		return nil, err
	}

	_, err = qtx.CreateTransaction(ctx, store.CreateTransactionParams{
		AccountID:   newAccount.ID,
		AmountCents: accountParams.BalanceCents,
		CategoryID:  1,
		Title:       "Initial balance",
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &newAccount, nil
}

func (s *AccountService) GetByID(accountID, userID int32) (*store.Account, error) {
	if accountID < 1 || userID < 1 {
		return nil, database.ErrRecordNotFound
	}

	account, err := s.queries.GetAccountByID(context.Background(), store.GetAccountByIDParams{
		ID:     accountID,
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

	return &account, nil
}

func (s *AccountService) DeleteByID(accountID, userID int32) error {
	if accountID < 1 || userID < 1 {
		return database.ErrRecordNotFound
	}

	result, err := s.queries.DeleteAccountById(context.Background(), store.DeleteAccountByIdParams{
		ID:     accountID,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrRecordNotFound
	}

	return nil
}

func (s *AccountService) GetSumBalance(accountID int32, userID int32) (int64, error) {
	if accountID < 1 || userID < 1 {
		return 0, database.ErrRecordNotFound
	}

	balance, err := s.queries.GetAccountSumBalance(context.Background(), store.GetAccountSumBalanceParams{
		AccountID: accountID,
		UserID:    userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, database.ErrRecordNotFound
		default:
			return 0, err
		}
	}

	return balance.Balance, nil
}

type UpdateAccountParams struct {
	Name *string            `json:"name"`
	Type *store.AccountType `json:"type"`
}

func (s *AccountService) UpdateByID(accountID, userID int32, updateParams *UpdateAccountParams) (*store.Account, error) {
	if accountID < 1 || userID < 1 {
		return nil, database.ErrRecordNotFound
	}

	ctx := context.Background()

	account, err := s.queries.GetAccountByID(ctx, store.GetAccountByIDParams{
		ID:     accountID,
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

	if updateParams.Name != nil {
		account.Name = *updateParams.Name
	}
	if updateParams.Type != nil {
		account.Type = *updateParams.Type
	}

	v := validator.New()
	if validateAccount(v, &account); !v.Valid() {
		return nil, v.GetErrors()
	}

	result, err := s.queries.UpdateAccountById(ctx, store.UpdateAccountByIdParams{
		Name:    account.Name,
		Type:    account.Type,
		ID:      account.ID,
		Version: account.Version,
		UserID:  account.UserID,
	})
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, database.ErrEditConflict
	}

	return &account, nil
}
