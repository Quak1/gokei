package testutils

import (
	"testing"

	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/internal/service"
)

func CreateTestUser(t *testing.T, userSvc *service.UserService, username string) *store.User {
	t.Helper()

	user, err := userSvc.Create(&service.InputUser{
		Username: username,
		Name:     "Test User",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return user
}
