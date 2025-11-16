package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/internal/testutils"
	"github.com/Quak1/gokei/pkg/assert"
)

func setupTestTransactionHandler(t *testing.T) (*TransactionHandler, *service.Service, func()) {
	db, cleanup, err := testutils.NewTestDB()
	if err != nil {
		cleanup()
		t.Fatal(err)
	}

	svc := service.New(db)
	handler := NewTransactionHandler(svc.Transaction)

	return handler, svc, cleanup
}

func TestTransactionHandler_Create(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestTransactionHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")
	account := testutils.CreateTestAccount(t, svc.Account, user.ID)
	category := testutils.CreateTestCategory(t, svc.Category)
	route := "/v1/transactions"

	tests := []struct {
		name           string
		requestBody    any
		expectedStatus int
		setupRequest   func(*testing.T) *http.Request
		validate       func(*testing.T, *http.Response)
	}{
		{
			name: "Create transaction",
			requestBody: map[string]any{
				"title":        "Test Transaction",
				"amount_cents": 10000,
				"account_id":   account.ID,
				"category_id":  category.ID,
			},
			expectedStatus: http.StatusCreated,
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*store.Transaction
				json.NewDecoder(r.Body).Decode(&resBody)

				transaction := resBody["transaction"]
				assert.Equal(t, transaction.Title, "Test Transaction")
				assert.Equal(t, transaction.AmountCents, 10000)
				assert.Equal(t, transaction.AccountID, user.ID)

				location := r.Header.Get("Location")
				assert.Equal(t, location, fmt.Sprintf("%s/%d", route, transaction.ID))
			},
		},
		{
			name: "Validation error",
			requestBody: map[string]any{
				"amount_cents": 10000,
				"account_id":   account.ID,
				"category_id":  category.ID,
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Category doesn't exist",
			requestBody: map[string]any{
				"title":        "Test Transaction",
				"amount_cents": 10000,
				"account_id":   account.ID,
				"category_id":  10,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Can't use initial category",
			requestBody: map[string]any{
				"title":        "Test Transaction",
				"amount_cents": 10000,
				"account_id":   account.ID,
				"category_id":  1,
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Invalid category ID",
			requestBody: map[string]any{
				"title":        "Test Transaction",
				"amount_cents": 10000,
				"account_id":   account.ID,
				"category_id":  -1,
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "User doesn't exist",
			requestBody: map[string]any{
				"title":        "Test Transaction",
				"amount_cents": 10000,
				"account_id":   10,
				"category_id":  category.ID,
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Fail to create transaction in another user's account",
			setupRequest: func(t *testing.T) *http.Request {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				account2 := testutils.CreateTestAccount(t, svc.Account, user2.ID)
				requestBody := map[string]any{
					"title":        "Test Transaction",
					"amount_cents": 10000,
					"account_id":   account2.ID,
					"category_id":  category.ID,
				}
				return testutils.CreatePostRequest(t, route, requestBody, user)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Bad JSON",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.setupRequest != nil {
				req = tt.setupRequest(t)
			} else {
				req = testutils.CreatePostRequest(t, route, tt.requestBody, user)
			}

			rr := httptest.NewRecorder()
			handler.Create(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, rs)
			}
		})
	}
}
