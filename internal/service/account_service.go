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
	queries *store.Queries
}

func NewAccountService(queries *store.Queries) *AccountService {
	return &AccountService{
		queries: queries,
	}
}

func validateAccount(v *validator.Validator, account *store.CreateAccountParams) {
	v.Check(validator.NonZero(account.Name), "name", "Must be provided")
	v.Check(validator.MaxLength(account.Name, 50), "name", "Must not be more than 50 bytes long")

	v.Check(validator.NonZero(account.Type), "type", "Must be provided")
	v.Check(validator.PermittedValue(account.Type, "debit", "cash", "credit"), "type", "Invalid account type. Valid types are credit, debit, and cash")
}

func (s *AccountService) GetAll() ([]*store.Account, error) {
	data, err := s.queries.GetAllAccounts(context.Background())
	if err != nil {
		return nil, err
	}

	accounts := make([]*store.Account, len(data))
	for i, a := range data {
		accounts[i] = &a
	}

	return accounts, nil
}

func (s *AccountService) Create(account *store.CreateAccountParams) (*store.Account, error) {
	v := validator.New()
	if validateAccount(v, account); !v.Valid() {
		return nil, v.GetErrors()
	}

	data, err := s.queries.CreateAccount(context.Background(), *account)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *AccountService) GetByID(id int32) (*store.Account, error) {
	if id < 1 {
		return nil, database.ErrRecordNotFound
	}

	account, err := s.queries.GetAccountByID(context.Background(), id)
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

func (s *AccountService) DeleteByID(id int32) error {
	if id < 1 {
		return database.ErrRecordNotFound
	}

	result, err := s.queries.DeleteAccountById(context.Background(), id)
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
