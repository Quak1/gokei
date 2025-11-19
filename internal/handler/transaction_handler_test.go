package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

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

func TestTransactionHandler_GetAll(t *testing.T) {
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
		expectedStatus int
		setup          func(*testing.T)
		validate       func(*testing.T, *http.Response)
	}{
		{
			name:           "Get all transactions",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) {
				testutils.CreateTestTransaction(t, svc.Transaction, user.ID, account.ID, category.ID)
			},
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*[]store.Transaction
				json.NewDecoder(r.Body).Decode(&resBody)

				transactions := resBody["transactions"]
				assert.Equal(t, len(*transactions), 2)
			},
		},
		{
			name:           "Only get user transactions",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				account2 := testutils.CreateTestAccount(t, svc.Account, user2.ID)
				testutils.CreateTestTransaction(t, svc.Transaction, user2.ID, account2.ID, category.ID)
				testutils.CreateTestTransaction(t, svc.Transaction, user2.ID, account2.ID, category.ID)
			},
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*[]store.Transaction
				json.NewDecoder(r.Body).Decode(&resBody)

				transactions := resBody["transactions"]
				assert.Equal(t, len(*transactions), 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}

			req := testutils.CreateGetRequest(t, route, user)

			rr := httptest.NewRecorder()
			handler.GetAll(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, res)
			}
		})
	}
}

func TestTransactionHandler_GetAccountTransactions(t *testing.T) {
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
		expectedStatus int
		accountID      string
		setup          func(*testing.T)
		validate       func(*testing.T, *http.Response)
	}{
		{
			name:           "Get all account transactions",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) {
				account2 := testutils.CreateTestAccount(t, svc.Account, user.ID)
				testutils.CreateTestTransaction(t, svc.Transaction, user.ID, account2.ID, category.ID)
			},
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*[]store.Transaction
				json.NewDecoder(r.Body).Decode(&resBody)

				transactions := resBody["transactions"]
				assert.Equal(t, len(*transactions), 1)
			},
		},
		{
			name:           "Only get user transactions",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) {
				testutils.CreateTestTransaction(t, svc.Transaction, user.ID, account.ID, category.ID)

				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				account2 := testutils.CreateTestAccount(t, svc.Account, user2.ID)
				testutils.CreateTestTransaction(t, svc.Transaction, user2.ID, account2.ID, category.ID)
				account3 := testutils.CreateTestAccount(t, svc.Account, user2.ID)
				testutils.CreateTestTransaction(t, svc.Transaction, user2.ID, account3.ID, category.ID)
			},
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*[]store.Transaction
				json.NewDecoder(r.Body).Decode(&resBody)

				transactions := resBody["transactions"]
				assert.Equal(t, len(*transactions), 2)
			},
		},
		{
			name:           "Not found account",
			expectedStatus: http.StatusNotFound,
			accountID:      "99",
		},
		{
			name:           "Invalid account ID",
			expectedStatus: http.StatusBadRequest,
			accountID:      "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}

			req := testutils.CreateGetRequest(t, route, user)
			if tt.accountID != "" {
				req.SetPathValue("accountID", tt.accountID)
			} else {
				req.SetPathValue("accountID", strconv.Itoa(int(account.ID)))
			}

			rr := httptest.NewRecorder()
			handler.GetAccountTransactions(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, res)
			}
		})
	}
}

func TestTransactionHandler_GetByID(t *testing.T) {
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
		transactionID  string
		expectedStatus int
		setup          func(*testing.T) int32
		validate       func(*testing.T, *http.Response)
	}{
		{
			name:           "Get transaction",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) int32 {
				transaction := testutils.CreateTestTransaction(t, svc.Transaction, user.ID, account.ID, category.ID)
				return transaction.ID
			},
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*store.Transaction
				json.NewDecoder(r.Body).Decode(&resBody)

				transaction := resBody["transaction"]
				assert.Equal(t, transaction.AccountID, account.ID)
				assert.Equal(t, transaction.CategoryID, category.ID)
				assert.Equal(t, transaction.Title, "Test Transaction")
			},
		},
		{
			name:           "Fail to get other user's transaction",
			expectedStatus: http.StatusNotFound,
			setup: func(t *testing.T) int32 {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				account2 := testutils.CreateTestAccount(t, svc.Account, user2.ID)
				transaction := testutils.CreateTestTransaction(t, svc.Transaction, user2.ID, account2.ID, category.ID)
				return transaction.ID
			},
		},
		{
			name:           "Not found",
			expectedStatus: http.StatusNotFound,
			transactionID:  "999",
		},
		{
			name:           "Negative ID",
			expectedStatus: http.StatusNotFound,
			transactionID:  "-1",
		},
		{
			name:           "Empty ID",
			expectedStatus: http.StatusBadRequest,
			transactionID:  "",
		},
		{
			name:           "Invalid ID",
			expectedStatus: http.StatusBadRequest,
			transactionID:  "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				id := tt.setup(t)
				tt.transactionID = strconv.Itoa(int(id))
			}

			req := testutils.CreateGetRequest(t, route, user)
			req.SetPathValue("transactionID", tt.transactionID)

			rr := httptest.NewRecorder()
			handler.GetByID(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, res)
			}
		})
	}
}

func TestTransactionHandler_DeleteByID(t *testing.T) {
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
		transactionID  string
		expectedStatus int
		setup          func(*testing.T) int32
		validate       func(*testing.T, *http.Response)
	}{
		{
			name:           "Delete transaction",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) int32 {
				transaction := testutils.CreateTestTransaction(t, svc.Transaction, user.ID, account.ID, category.ID)

				transactions, err := svc.Transaction.GetAll(user.ID)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, len(transactions), 2)

				return transaction.ID
			},
			validate: func(t *testing.T, r *http.Response) {
				transactions, err := svc.Transaction.GetAll(user.ID)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, len(transactions), 1)
			},
		},
		{
			name:           "Fail to delete initial transaction",
			expectedStatus: http.StatusForbidden,
			setup: func(t *testing.T) int32 {
				transactions, err := svc.Transaction.GetAll(user.ID)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, len(transactions), 1)

				return transactions[0].ID
			},
		},
		{
			name:           "Fail to delete other user's transaction",
			expectedStatus: http.StatusNotFound,
			setup: func(t *testing.T) int32 {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				account2 := testutils.CreateTestAccount(t, svc.Account, user2.ID)
				transaction := testutils.CreateTestTransaction(t, svc.Transaction, user2.ID, account2.ID, category.ID)
				return transaction.ID
			},
		},
		{
			name:           "Not found",
			expectedStatus: http.StatusNotFound,
			transactionID:  "999",
		},
		{
			name:           "Negative ID",
			expectedStatus: http.StatusNotFound,
			transactionID:  "-1",
		},
		{
			name:           "Empty ID",
			expectedStatus: http.StatusBadRequest,
			transactionID:  "",
		},
		{
			name:           "Invalid ID",
			expectedStatus: http.StatusBadRequest,
			transactionID:  "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				id := tt.setup(t)
				tt.transactionID = strconv.Itoa(int(id))
			}

			req := testutils.CreateGetRequest(t, route, user)
			req.SetPathValue("transactionID", tt.transactionID)

			rr := httptest.NewRecorder()
			handler.DeleteByID(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, res)
			}
		})
	}
}

func TestTransactionHandler_UpdateByID(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestTransactionHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")
	account := testutils.CreateTestAccount(t, svc.Account, user.ID)
	category := testutils.CreateTestCategory(t, svc.Category)
	transaction := testutils.CreateTestTransaction(t, svc.Transaction, user.ID, account.ID, category.ID)
	transactionID := strconv.Itoa(int(transaction.ID))
	route := "/v1/transactions"

	tests := []struct {
		name           string
		body           any
		transactionID  string
		expectedStatus int
		setup          func(*testing.T) int32
		validate       func(*testing.T, *http.Response)
	}{
		{
			name: "Update transaction",
			body: map[string]any{
				"title":        "updated transaction",
				"amount_cents": 321,
			},
			transactionID:  transactionID,
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*store.Transaction
				json.NewDecoder(r.Body).Decode(&resBody)

				resTransaction := resBody["transaction"]
				transaction, err := svc.Transaction.GetByID(resTransaction.ID, user.ID)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, transaction.Title, "updated transaction")
				assert.Equal(t, transaction.AmountCents, 321)
				assert.Equal(t, transaction.CategoryID, category.ID)
				assert.Equal(t, transaction.Note, "")
				assert.Equal(t, transaction.Attachment, "")
			},
		},
		{
			name: "Update full transaction",
			body: map[string]any{
				"amount_cents": 1234,
				"account_id":   2,
				"category_id":  3,
				"title":        "Updated transaction title",
				"date":         "2025-11-17T00:00:00Z",
				"attachment":   "att",
				"note":         "updated transaction note",
			},
			transactionID:  transactionID,
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) int32 {
				testutils.CreateTestAccount(t, svc.Account, user.ID)
				testutils.CreateTestCategory(t, svc.Category)
				return transaction.ID
			},
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*store.Transaction
				json.NewDecoder(r.Body).Decode(&resBody)

				resTransaction := resBody["transaction"]
				transaction, err := svc.Transaction.GetByID(resTransaction.ID, user.ID)
				if err != nil {
					t.Fatal(err)
				}

				date, err := time.Parse(time.DateOnly, "2025-11-17")
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, transaction.AmountCents, 1234)
				assert.Equal(t, transaction.AccountID, 2)
				assert.Equal(t, transaction.CategoryID, 3)
				assert.Equal(t, transaction.Title, "Updated transaction title")
				assert.Equal(t, transaction.Date.UTC(), date)
				assert.Equal(t, transaction.Attachment, "att")
				assert.Equal(t, transaction.Note, "updated transaction note")
			},
		},
		{
			name: "Fail to update other user's transaction",
			body: map[string]any{
				"title": "new title",
			},
			expectedStatus: http.StatusNotFound,
			setup: func(t *testing.T) int32 {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				account2 := testutils.CreateTestAccount(t, svc.Account, user2.ID)
				transaction := testutils.CreateTestTransaction(t, svc.Transaction, user2.ID, account2.ID, category.ID)
				return transaction.ID
			},
		},
		{
			name:          "Can't use initial category",
			transactionID: transactionID,
			body: map[string]any{
				"title":       "New title",
				"category_id": 1,
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Invalid category",
			body: map[string]any{
				"category_id": 99,
			},
			transactionID:  transactionID,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid key",
			body: map[string]any{
				"test": "test",
			},
			transactionID:  transactionID,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Validation error",
			body: map[string]any{
				"title": "",
			},
			transactionID:  transactionID,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Incorrect JSON",
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Not found",
			expectedStatus: http.StatusNotFound,
			transactionID:  "999",
		},
		{
			name:           "Negative ID",
			expectedStatus: http.StatusNotFound,
			transactionID:  "-1",
		},
		{
			name:           "Empty ID",
			expectedStatus: http.StatusBadRequest,
			transactionID:  "",
		},
		{
			name:           "Invalid ID",
			expectedStatus: http.StatusBadRequest,
			transactionID:  "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				id := tt.setup(t)
				tt.transactionID = strconv.Itoa(int(id))
			}

			req := testutils.CreatePostRequest(t, route, tt.body, user)
			req.SetPathValue("transactionID", tt.transactionID)

			rr := httptest.NewRecorder()
			handler.UpdateByID(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, res)
			}
		})
	}
}
