package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Quak1/gokei/internal/appcontext"
	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/internal/testutils"
	"github.com/Quak1/gokei/pkg/assert"
)

func setupTestUserHandler(t *testing.T) (*UserHandler, *service.Service, func()) {
	db, cleanup, err := testutils.NewTestDB()
	if err != nil {
		cleanup()
		t.Fatal(err)
	}

	svc := service.New(db)
	handler := NewUserHandler(svc.User)

	return handler, svc, cleanup
}

func TestUserHandler_Create(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestUserHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")
	route := "/v1/users"

	tests := []struct {
		name           string
		requestBody    any
		expectedStatus int
		validate       func(*testing.T, *http.Response)
	}{
		{
			name: "Create user",
			requestBody: map[string]any{
				"name":     "Test User",
				"username": "username",
				"password": "TestPassword",
			},
			expectedStatus: http.StatusCreated,
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*store.User
				json.NewDecoder(r.Body).Decode(&resBody)

				user := resBody["user"]
				assert.Equal(t, user.Name, "Test User")
				assert.Equal(t, user.Username, "username")

				location := r.Header.Get("Location")
				assert.Equal(t, location, fmt.Sprintf("%s/%d", route, user.ID))
			},
		},
		{
			name: "Fail to create user with duplicate username",
			requestBody: map[string]any{
				"name":     "Test User",
				"username": "username",
				"password": "TestPassword",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Validation error - missing username",
			requestBody: map[string]any{
				"name":     "Test User",
				"password": "TestPassword",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Bad JSON",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := testutils.CreatePostRequest(t, route, tt.requestBody, user)

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

func TestUserHandler_GetByID(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestUserHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")
	userID := strconv.Itoa(int(user.ID))
	route := "/v1/users"

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		setup          func(*testing.T) *http.Request
		validate       func(*testing.T, *http.Response)
	}{
		{
			name:           "Get user",
			expectedStatus: http.StatusOK,
			userID:         userID,
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*store.User
				json.NewDecoder(r.Body).Decode(&resBody)

				resUser := resBody["user"]
				assert.Equal(t, resUser.ID, user.ID)
				assert.Equal(t, resUser.Username, "testuser")
			},
		},
		{
			name:           "Fail to get other user's details",
			expectedStatus: http.StatusNotFound,
			setup: func(t *testing.T) *http.Request {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")

				req := testutils.CreateGetRequest(t, route, user)
				req.SetPathValue("userID", strconv.Itoa(int(user2.ID)))

				return req
			},
		},
		{
			name:           "Fail to get details",
			expectedStatus: http.StatusNotFound,
			setup: func(t *testing.T) *http.Request {
				user2 := testutils.CreateTestUser(t, svc.User, "deleteduser")
				err := svc.User.DeleteByID(user2.ID)
				if err != nil {
					t.Fatal(err)
				}

				req := testutils.CreateGetRequest(t, route, user2)
				req.SetPathValue("userID", strconv.Itoa(int(user2.ID)))

				return req
			},
		},
		{
			name:           "Not found",
			expectedStatus: http.StatusNotFound,
			userID:         "999",
		},
		{
			name:           "Negative ID",
			expectedStatus: http.StatusNotFound,
			userID:         "-1",
		},
		{
			name:           "Empty ID",
			expectedStatus: http.StatusBadRequest,
			userID:         "",
		},
		{
			name:           "Invalid ID",
			expectedStatus: http.StatusBadRequest,
			userID:         "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request

			if tt.setup != nil {
				req = tt.setup(t)
			} else {
				req = testutils.CreateGetRequest(t, route, user)
				req.SetPathValue("userID", tt.userID)
			}

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

func TestUserHandler_DeleteByID(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestUserHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")
	route := "/v1/users"

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		setup          func(*testing.T) (*http.Request, string)
		validate       func(*testing.T, *http.Response, int32)
	}{
		{
			name:           "Delete user",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) (*http.Request, string) {
				user := testutils.CreateTestUser(t, svc.User, "deleteuser")
				userID := strconv.Itoa(int(user.ID))

				req := httptest.NewRequest(http.MethodDelete, route, nil)
				req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
					ID:       user.ID,
					Username: user.Username,
				})
				req.SetPathValue("userID", userID)

				return req, userID
			},
			validate: func(t *testing.T, r *http.Response, userID int32) {
				_, err := svc.User.GetByID(userID)
				if !errors.Is(err, database.ErrRecordNotFound) {
					t.Error("Expected user to be deleted, but found it")
				}
			},
		},
		{
			name:           "Fail to delete other user",
			expectedStatus: http.StatusNotFound,
			setup: func(t *testing.T) (*http.Request, string) {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				user2ID := strconv.Itoa(int(user2.ID))

				req := httptest.NewRequest(http.MethodDelete, route, nil)
				req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
					ID:       user.ID,
					Username: user.Username,
				})
				req.SetPathValue("userID", user2ID)

				return req, user2ID
			},
			validate: func(t *testing.T, r *http.Response, userID int32) {
				user, err := svc.User.GetByID(userID)
				if err != nil {
					t.Fatal("Failed to get user 2")
				}

				assert.Equal(t, user.Username, "user2")
			},
		},
		{
			name:           "Fail to delete deleted user",
			expectedStatus: http.StatusNotFound,
			setup: func(t *testing.T) (*http.Request, string) {
				user2 := testutils.CreateTestUser(t, svc.User, "deleteduser")
				user2ID := strconv.Itoa(int(user2.ID))

				err := svc.User.DeleteByID(user2.ID)
				if err != nil {
					t.Fatal(err)
				}

				req := httptest.NewRequest(http.MethodDelete, route, nil)
				req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
					ID:       user2.ID,
					Username: user2.Username,
				})
				req.SetPathValue("userID", user2ID)

				return req, user2ID
			},
		},
		{
			name:           "Not found",
			expectedStatus: http.StatusNotFound,
			userID:         "999",
		},
		{
			name:           "Negative ID",
			expectedStatus: http.StatusNotFound,
			userID:         "-1",
		},
		{
			name:           "Empty ID",
			expectedStatus: http.StatusBadRequest,
			userID:         "",
		},
		{
			name:           "Invalid ID",
			expectedStatus: http.StatusBadRequest,
			userID:         "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request

			if tt.setup != nil {
				req, tt.userID = tt.setup(t)
			} else {
				req = httptest.NewRequest(http.MethodDelete, route, nil)
				req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
					ID:       user.ID,
					Username: user.Username,
				})
				req.SetPathValue("userID", tt.userID)
			}

			rr := httptest.NewRecorder()
			handler.DeleteByID(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				id, err := strconv.ParseInt(tt.userID, 10, 32)
				if err != nil {
					t.Fatal(err)
				}

				tt.validate(t, res, int32(id))
			}
		})
	}
}

func TestUserHandler_UpdateByID(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestUserHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")
	userID := strconv.Itoa(int(user.ID))
	route := "/v1/users"

	tests := []struct {
		name           string
		body           any
		userID         string
		expectedStatus int
		setup          func(*testing.T) *http.Request
		validate       func(*testing.T, *http.Response)
	}{
		{
			name: "Update User",
			body: map[string]any{
				"name":     "New user name",
				"password": "newpassword",
			},
			userID:         userID,
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*store.User
				json.NewDecoder(r.Body).Decode(&resBody)

				resUser := resBody["user"]
				user, err := svc.User.GetByID(resUser.ID)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, user.Name, "New user name")
				assert.Equal(t, user.Version, 2)
			},
		},
		{
			name:           "Fail to update deleted user",
			expectedStatus: http.StatusNotFound,
			setup: func(t *testing.T) *http.Request {
				user2 := testutils.CreateTestUser(t, svc.User, "deleteduser")
				body := map[string]any{
					"name": "new name",
				}

				err := svc.User.DeleteByID(user2.ID)
				if err != nil {
					t.Fatal(err)
				}

				req := testutils.CreatePostRequest(t, route, body, user2)
				req.SetPathValue("userID", strconv.Itoa(int(user2.ID)))

				return req
			},
		},
		{
			name:           "Fail to update other user's details",
			expectedStatus: http.StatusNotFound,
			setup: func(t *testing.T) *http.Request {
				user2 := testutils.CreateTestUser(t, svc.User, "user2")
				body := map[string]any{
					"name": "new name",
				}

				req := testutils.CreatePostRequest(t, route, body, user2)
				req.SetPathValue("userID", userID)

				return req
			},
		},
		{
			name: "Fail to update username",
			body: map[string]any{
				"username": "new_username",
			},
			userID:         userID,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Validation error",
			body: map[string]any{
				"password": "",
			},
			userID:         userID,
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
			userID:         "999",
		},
		{
			name:           "Negative ID",
			expectedStatus: http.StatusNotFound,
			userID:         "-1",
		},
		{
			name:           "Empty ID",
			expectedStatus: http.StatusBadRequest,
			userID:         "",
		},
		{
			name:           "Invalid ID",
			expectedStatus: http.StatusBadRequest,
			userID:         "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.setup != nil {
				req = tt.setup(t)
			} else {
				req = testutils.CreatePostRequest(t, route, tt.body, user)
				req.SetPathValue("userID", tt.userID)
			}

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
