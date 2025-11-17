package testutils

import (
	"testing"

	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/internal/service"
)

func CreateTestTransaction(t *testing.T, svc *service.TransactionService, userID, accountID, categoryID int32) *store.Transaction {
	t.Helper()

	transaction, err := svc.Create(userID, &store.CreateTransactionParams{
		Title:       "Test Transaction",
		AccountID:   accountID,
		AmountCents: 100000,
		CategoryID:  categoryID,
	})
	if err != nil {
		t.Fatalf("failed to create test transaction: %v", err)
	}

	return transaction
}
