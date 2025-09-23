package services

import (
	"context"
	"strconv"

	"github.com/Quak1/gokei/internal/apperrors"
	"github.com/Quak1/gokei/internal/database/queries"
)

type AccountService struct {
	queries *queries.Queries
}

func NewAccountService(queries *queries.Queries) *AccountService {
	return &AccountService{
		queries: queries,
	}
}

type CreateAccountRequest struct {
	Name string              `json:"name" validate:"required,min=2"`
	Kind queries.AccountType `json:"accountType" validate:"required,oneof=credit debit cash"`
}

func (s *AccountService) CreateAccount(ctx context.Context, userID string, req *CreateAccountRequest) (*queries.Account, error) {
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeBadRequest, "Invalid user ID", err)
	}

	account, err := s.queries.CreateAccount(ctx, queries.CreateAccountParams{
		OwnerID: int32(userIDInt),
		Kind:    req.Kind,
		Name:    req.Name,
	})
	if err != nil {
		return nil, apperrors.New(apperrors.CodeInternal, "Failed to create account", err)
	}

	return &account, nil
}
