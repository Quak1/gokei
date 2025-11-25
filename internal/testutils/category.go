package testutils

import (
	"testing"

	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/internal/service"
)

func CreateTestCategory(t *testing.T, svc *service.CategoryService, userID int32) *store.Category {
	t.Helper()

	category, err := svc.Create(&store.CreateCategoryParams{
		Name:   "Test Category",
		Color:  "#A1B2C3",
		Icon:   "T",
		UserID: userID,
	})
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	return category
}
