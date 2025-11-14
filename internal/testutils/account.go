package testutils

import (
	"testing"

	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/internal/service"
)

func CreateTestAccount(t *testing.T, accountSvc *service.AccountService, userID int32, name ...string) *store.Account {
	t.Helper()

	accountName := "Test account"
	if len(name) > 0 {
		accountName = name[0]
	}

	account, err := accountSvc.Create(&store.CreateAccountParams{
		Type:         store.AccountTypeDebit,
		Name:         accountName,
		UserID:       userID,
		BalanceCents: 10000,
	})
	if err != nil {
		t.Fatalf("failed to create test account: %v", err)
	}

	return account
}
