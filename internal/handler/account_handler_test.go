package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Quak1/gokei/internal/appcontext"
	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/internal/testutils"
	"github.com/Quak1/gokei/pkg/assert"
)

func setupTestAccountHandler(t *testing.T) (*AccountHandler, *service.Service, func()) {
	db, cleanup, err := testutils.NewTestDB()
	if err != nil {
		t.Fatalf("test db setup failed: %v", err)
	}

	svc := service.New(db)
	handler := NewAccountHandler(svc.Account)

	return handler, svc, cleanup
}

func TestAccountHandler_Create(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestAccountHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")

	route := "/v1/accounts"

	tests := []struct {
		name           string
		requestBody    any
		expectedStatus int
		validate       func(*testing.T, *http.Response)
	}{
		{
			name: "Create account",
			requestBody: map[string]any{
				"name":            "Test Account",
				"type":            "credit",
				"initial_balance": 10000,
			},
			expectedStatus: http.StatusCreated,
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*store.Account
				json.NewDecoder(r.Body).Decode(&resBody)

				account := resBody["account"]
				assert.Equal(t, account.Name, "Test Account")
				assert.Equal(t, account.Type, "credit")
				assert.Equal(t, account.BalanceCents, 10000)

				location := r.Header.Get("Location")
				assert.Equal(t, location, fmt.Sprintf("%s/%d", route, account.ID))
			},
		},
		{
			name:           "Incorrect JSON",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Validation error",
			requestBody: map[string]any{
				"type": "red",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "validation error - invalid type",
			expectedStatus: http.StatusUnprocessableEntity,
			requestBody: map[string]any{
				"name":            "Test Account",
				"type":            "invalid",
				"initial_balance": 10000,
			},
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]any
				json.NewDecoder(r.Body).Decode(&resBody)

				errors := resBody["error"].(map[string]any)
				_, hasNameError := errors["type"]
				assert.Equal(t, hasNameError, true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
			req.Header.Set("Contenty-Type", "application/json")
			req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
				ID:       user.ID,
				Username: user.Username,
			})

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

func TestAccountHandler_GetAll(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestAccountHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")

	route := "/v1/accounts"

	tests := []struct {
		name           string
		expectedStatus int
		setup          func(*testing.T)
		validate       func(*testing.T, *http.Response)
	}{
		{
			name:           "Get None",
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, rs *http.Response) {
				var resBody map[string][]*store.Account
				json.NewDecoder(rs.Body).Decode(&resBody)

				accounts := resBody["accounts"]
				assert.Equal(t, len(accounts), 0)
			},
		},
		{
			name:           "Only get user accounts",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				testutils.CreateTestAccount(t, svc.Account, user2.ID)
				testutils.CreateTestAccount(t, svc.Account, user2.ID)

				testutils.CreateTestAccount(t, svc.Account, user.ID)
			},
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string][]*store.Account
				json.NewDecoder(r.Body).Decode(&resBody)

				accounts := resBody["accounts"]
				assert.Equal(t, len(accounts), 1)
				assert.Equal(t, accounts[0].UserID, user.ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			}

			req := httptest.NewRequest(http.MethodPost, route, nil)
			req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
				ID:       user.ID,
				Username: user.Username,
			})

			rr := httptest.NewRecorder()
			handler.GetAll(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, rs)
			}
		})
	}
}

func TestAccountHandler_GetByID(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestAccountHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")

	route := "/v1/accounts"
	idPath := "accountID"

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		setup          func(*testing.T) int32
		validate       func(*testing.T, *http.Response)
	}{
		{
			name:           "Get account",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) int32 {
				testutils.CreateTestAccount(t, svc.Account, user.ID)
				account := testutils.CreateTestAccount(t, svc.Account, user.ID, "found account")
				return account.ID
			},
			validate: func(t *testing.T, rs *http.Response) {
				var resBody map[string]*store.Account
				json.NewDecoder(rs.Body).Decode(&resBody)

				account := resBody["account"]
				assert.Equal(t, account.Name, "found account")
			},
		},
		{
			name:           "Fail to get other user's account",
			expectedStatus: http.StatusNotFound,
			setup: func(t *testing.T) int32 {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				testutils.CreateTestAccount(t, svc.Account, user.ID)
				account := testutils.CreateTestAccount(t, svc.Account, user2.ID)
				return account.ID
			},
		},
		{
			name:           "Not found",
			expectedStatus: http.StatusNotFound,
			id:             "999",
		},
		{
			name:           "Negative ID",
			expectedStatus: http.StatusNotFound,
			id:             "-5",
		},
		{
			name:           "Empty ID",
			expectedStatus: http.StatusBadRequest,
			id:             "",
		},
		{
			name:           "Invalid ID",
			expectedStatus: http.StatusBadRequest,
			id:             "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				id := tt.setup(t)
				tt.id = strconv.Itoa(int(id))
			}

			req := httptest.NewRequest(http.MethodPost, route, nil)
			req.SetPathValue(idPath, tt.id)
			req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
				ID:       user.ID,
				Username: user.Username,
			})

			rr := httptest.NewRecorder()
			handler.GetByID(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, rs)
			}
		})
	}
}

func TestAccountHandler_DeleteByID(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestAccountHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")

	route := "/v1/accounts"
	idPath := "accountID"

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		setup          func(*testing.T) int32
		validate       func(*testing.T, *http.Response)
	}{
		{
			name:           "Delete account",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) int32 {
				account := testutils.CreateTestAccount(t, svc.Account, user.ID)
				return account.ID
			},
			validate: func(t *testing.T, rs *http.Response) {
				accounts, err := svc.Account.GetAll(user.ID)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, len(accounts), 0)
			},
		},
		{
			name:           "Fail to delete other users account",
			expectedStatus: http.StatusNotFound,
			setup: func(t *testing.T) int32 {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				account := testutils.CreateTestAccount(t, svc.Account, user2.ID)
				return account.ID
			},
		},
		{
			name:           "Not found",
			expectedStatus: http.StatusNotFound,
			id:             "999",
		},
		{
			name:           "Negative ID",
			expectedStatus: http.StatusNotFound,
			id:             "-5",
		},
		{
			name:           "Empty ID",
			expectedStatus: http.StatusBadRequest,
			id:             "",
		},
		{
			name:           "Invalid ID",
			expectedStatus: http.StatusBadRequest,
			id:             "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				id := tt.setup(t)
				tt.id = strconv.Itoa(int(id))
			}

			req := httptest.NewRequest(http.MethodPost, route, nil)
			req.SetPathValue(idPath, tt.id)
			req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
				ID:       user.ID,
				Username: user.Username,
			})

			rr := httptest.NewRecorder()
			handler.DeleteByID(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, rs)
			}
		})
	}
}

func TestAccountHandler_GetSumBalance(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestAccountHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")

	route := "/v1/accounts"
	idPath := "accountID"

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		setup          func(*testing.T) int32
		validate       func(*testing.T, *http.Response)
	}{
		{
			name:           "Get balance",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) int32 {
				account := testutils.CreateTestAccount(t, svc.Account, user.ID)

				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				testutils.CreateTestAccount(t, svc.Account, user2.ID)

				return account.ID
			},
			validate: func(t *testing.T, rs *http.Response) {
				var resBody map[string]int64
				json.NewDecoder(rs.Body).Decode(&resBody)

				balance := resBody["balance_cents"]
				assert.Equal(t, balance, 10000)
			},
		},
		{
			name:           "Not found",
			expectedStatus: http.StatusNotFound,
			id:             "999",
		},
		{
			name:           "Negative ID",
			expectedStatus: http.StatusNotFound,
			id:             "-5",
		},
		{
			name:           "Empty ID",
			expectedStatus: http.StatusBadRequest,
			id:             "",
		},
		{
			name:           "Invalid ID",
			expectedStatus: http.StatusBadRequest,
			id:             "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				id := tt.setup(t)
				tt.id = strconv.Itoa(int(id))
			}

			req := httptest.NewRequest(http.MethodPost, route, nil)
			req.SetPathValue(idPath, tt.id)
			req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
				ID:       user.ID,
				Username: user.Username,
			})

			rr := httptest.NewRecorder()
			handler.GetSumBalance(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, rs)
			}
		})
	}
}

func TestAccountHandler_UpdateByID(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestAccountHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")
	account := testutils.CreateTestAccount(t, svc.Account, user.ID)
	accountID := strconv.Itoa(int(account.ID))

	route := "/v1/accounts"
	idPath := "accountID"

	tests := []struct {
		name           string
		requestBody    any
		id             string
		expectedStatus int
		setup          func(*testing.T) int32
		validate       func(*testing.T, *http.Response)
	}{
		{
			name: "Update name",
			requestBody: map[string]any{
				"name": "Update name",
			},
			id:             accountID,
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, rs *http.Response) {
				var resBody map[string]*store.Account
				json.NewDecoder(rs.Body).Decode(&resBody)

				updatedAccount := resBody["account"]
				assert.Equal(t, updatedAccount.Name, "Update name")
				assert.Equal(t, updatedAccount.Type, account.Type)
				assert.Equal(t, updatedAccount.BalanceCents, account.BalanceCents)
			},
		},
		{
			name: "Fail to update other user's account",
			requestBody: map[string]any{
				"name": "Update name",
			},
			setup: func(t *testing.T) int32 {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				account := testutils.CreateTestAccount(t, svc.Account, user2.ID)
				return account.ID
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Validation error",
			requestBody: map[string]any{
				"type": "red",
			},
			id:             accountID,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "validation error - invalid key",
			id:             accountID,
			expectedStatus: http.StatusBadRequest,
			requestBody: map[string]any{
				"balance": 1000000,
			},
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]any
				json.NewDecoder(r.Body).Decode(&resBody)

				error := resBody["error"].(string)
				assert.StringContains(t, error, "balance")
			},
		},
		{
			name:           "Incorrect JSON",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Not found",
			requestBody:    map[string]any{},
			expectedStatus: http.StatusNotFound,
			id:             "999",
		},
		{
			name:           "Negative ID",
			requestBody:    map[string]any{},
			expectedStatus: http.StatusNotFound,
			id:             "-5",
		},
		{
			name:           "Empty ID",
			requestBody:    map[string]any{},
			expectedStatus: http.StatusBadRequest,
			id:             "",
		},
		{
			name:           "Invalid ID",
			requestBody:    map[string]any{},
			expectedStatus: http.StatusBadRequest,
			id:             "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				id := tt.setup(t)
				tt.id = strconv.Itoa(int(id))
			}

			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
			req.Header.Set("Contenty-Type", "application/json")
			req.SetPathValue(idPath, tt.id)
			req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
				ID:       user.ID,
				Username: user.Username,
			})

			rr := httptest.NewRecorder()
			handler.UpdateByID(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, rs)
			}
		})
	}
}

func TestAccountHandler_TransferByID(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestAccountHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")
	account := testutils.CreateTestAccount(t, svc.Account, user.ID)
	account2 := testutils.CreateTestAccount(t, svc.Account, user.ID)

	accountID := strconv.Itoa(int(account.ID))
	target := "/v1/accounts"
	idPath := "accountID"

	tests := []struct {
		name           string
		requestBody    any
		senderID       string
		expectedStatus int
		setup          func(*testing.T) *http.Request
		validate       func(*testing.T, *http.Response)
	}{
		{
			name: "Transfer funds",
			requestBody: map[string]any{
				"amount":       1000,
				"recipient_id": account2.ID,
			},
			senderID:       accountID,
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, rs *http.Response) {
				var resBody map[string]*store.Transaction
				json.NewDecoder(rs.Body).Decode(&resBody)

				transaction := resBody["transaction"]
				assert.Equal(t, transaction.AmountCents, -1000)
				assert.StringContains(t, transaction.Title, "TRANSFER")

				fetchAccount, err := svc.Account.GetByID(account.ID, user.ID)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, fetchAccount.BalanceCents, account.BalanceCents-1000)
			},
		},
		{
			name: "Fail to transfer from other user's account",
			setup: func(t *testing.T) *http.Request {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				user2Account := testutils.CreateTestAccount(t, svc.Account, user2.ID)

				requestBody := map[string]any{
					"amount":       1000,
					"recipient_id": user2Account.ID,
				}

				req := testutils.CreatePostRequest(t, target, requestBody, user2)
				req.SetPathValue(idPath, accountID)

				return req
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Validation error - fail to transfer to same account",
			requestBody: map[string]any{
				"amount":       1000,
				"recipient_id": account.ID,
			},
			senderID:       accountID,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Validation error - nagative amount",
			requestBody: map[string]any{
				"amount":       -100,
				"recipient_id": account2.ID,
			},
			senderID:       accountID,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Incorrect JSON",
			requestBody:    "",
			senderID:       accountID,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Not found",
			requestBody: map[string]any{
				"amount":       100,
				"recipient_id": account2.ID,
			},
			senderID:       "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Negative ID",
			requestBody: map[string]any{
				"amount":       100,
				"recipient_id": account2.ID,
			},
			senderID:       "-5",
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Empty ID",
			requestBody: map[string]any{
				"amount":       100,
				"recipient_id": account2.ID,
			},
			senderID:       "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid ID",
			requestBody: map[string]any{
				"amount":       100,
				"recipient_id": account2.ID,
			},
			senderID:       "test",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.setup != nil {
				req = tt.setup(t)
			} else {
				body, err := json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}

				req = httptest.NewRequest(http.MethodPost, target, bytes.NewBuffer(body))
				req.Header.Set("Contenty-Type", "application/json")
				req.SetPathValue(idPath, tt.senderID)
				req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
					ID:       user.ID,
					Username: user.Username,
				})
			}

			rr := httptest.NewRecorder()
			handler.TransferByID(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, rs)
			}
		})
	}
}
