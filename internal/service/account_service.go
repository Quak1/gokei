package service

import (
	"context"

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
