package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/internal/testutils"
	"github.com/Quak1/gokei/pkg/assert"
)

func setupTestCategoryHandler(t *testing.T) (*CategoryHandler, *service.Service, func()) {
	db, cleanup, err := testutils.NewTestDB()
	if err != nil {
		t.Fatalf("test db setup failed: %v", err)
	}

	svc := service.New(db)
	handler := NewCategoryHandler(svc.Category)

	return handler, svc, cleanup
}

func TestCategoryHandler_Create(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestCategoryHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")

	route := "/v1/categories"

	tests := []struct {
		name           string
		requestBody    any
		expectedStatus int
		setupRequest   func(*testing.T) *http.Request
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "Create category",
			requestBody: map[string]any{
				"name":  "Test Category",
				"color": "#FFF",
				"icon":  "T",
			},
			checkResponse: func(t *testing.T, rs *http.Response) {
				var resBody map[string]*store.Category
				json.NewDecoder(rs.Body).Decode(&resBody)

				category := resBody["category"]
				assert.Equal(t, category.Name, "Test Category")
				assert.Equal(t, category.Color, "#FFF")
				assert.Equal(t, category.Icon, "T")

				location := rs.Header.Get("Location")
				assert.Equal(t, location, fmt.Sprintf("%s/%d", route, category.ID))
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid user",
			setupRequest: func(t *testing.T) *http.Request {
				requestBody := map[string]any{
					"name":  "Test Category",
					"color": "#FFF",
					"icon":  "T",
				}

				req := testutils.CreatePostRequest(t, route, requestBody, &store.User{
					ID:       5,
					Username: "Fake",
				})

				return req
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Incorrect JSON",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Validation error",
			requestBody: map[string]any{
				"color": "FFF",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := testutils.CreatePostRequest(t, route, tt.requestBody, user)
			if tt.setupRequest != nil {
				req = tt.setupRequest(t)
			}

			rr := httptest.NewRecorder()
			handler.Create(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(t, rs)
			}
		})
	}
}

func TestCategoryHandler_GetAll(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestCategoryHandler(t)
	defer cleanup()

	route := "/v1/categories"
	user := testutils.CreateTestUser(t, svc.User, "testuser")

	tests := []struct {
		name           string
		expectedStatus int
		setup          func(*testing.T)
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "Get all categories",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rs *http.Response) {
				var resBody map[string][]*store.Category
				json.NewDecoder(rs.Body).Decode(&resBody)

				categories := resBody["categories"]
				assert.Equal(t, len(categories), 1)
			},
		},
		{
			name:           "Only get own categories",
			expectedStatus: http.StatusOK,
			setup: func(t *testing.T) {
				user2 := testutils.CreateTestUser(t, svc.User, "testuser")
				testutils.CreateTestCategory(t, svc.Category, user2.ID)
			},
			checkResponse: func(t *testing.T, rs *http.Response) {
				var resBody map[string][]*store.Category
				json.NewDecoder(rs.Body).Decode(&resBody)

				categories := resBody["categories"]
				assert.Equal(t, len(categories), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := testutils.CreateGetRequest(t, route, user)

			rr := httptest.NewRecorder()
			handler.GetAll(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(t, rs)
			}
		})
	}
}

func TestCategoryHandler_GetByID(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestCategoryHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")
	category := testutils.CreateTestCategory(t, svc.Category, user.ID)
	categoryID := strconv.Itoa(int(category.ID))

	route := "/v1/categories"
	idPath := "categoryID"

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		setupRequest   func(*testing.T) *http.Request
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "Get category",
			id:             categoryID,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Fail to get other user category",
			id:   categoryID,
			setupRequest: func(t *testing.T) *http.Request {
				user2 := testutils.CreateTestUser(t, svc.User, "testuser2")
				category2 := testutils.CreateTestCategory(t, svc.Category, user2.ID)

				req := testutils.CreateGetRequest(t, route, user)
				req.SetPathValue(idPath, strconv.Itoa(int(category2.ID)))

				return req
			},

			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Not found",
			id:             "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid ID",
			id:             "bad",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty ID",
			id:             "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Negative ID",
			id:             "-5",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := testutils.CreateGetRequest(t, route, user)
			req.SetPathValue(idPath, tt.id)

			if tt.setupRequest != nil {
				req = tt.setupRequest(t)
			}

			rr := httptest.NewRecorder()
			handler.GetByID(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(t, rs)
			}
		})
	}
}

func TestCategoryHandler_DeleteByID(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestCategoryHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")
	category, err := svc.Category.Create(&store.CreateCategoryParams{
		UserID: user.ID,
		Name:   "Test",
		Color:  "#123",
		Icon:   "T",
	})
	if err != nil {
		t.Fatal(err)
	}
	categoryID := strconv.Itoa(int(category.ID))

	route := "/v1/categories"
	idPath := "categoryID"

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		setup          func(*testing.T) *http.Request
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "Delete category",
			id:             categoryID,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Fail to delete other user's category",
			setup: func(t *testing.T) *http.Request {
				user2 := testutils.CreateTestUser(t, svc.User, "testuser2")
				category2 := testutils.CreateTestCategory(t, svc.Category, user2.ID)

				req := testutils.CreateGetRequest(t, route, user)
				req.SetPathValue(idPath, strconv.Itoa(int(category2.ID)))

				return req
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Fail to delete initial category",
			id:             "1",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Not found",
			id:             "2",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid ID",
			id:             "bad",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty ID",
			id:             "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Negative ID",
			id:             "-1",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := testutils.CreateGetRequest(t, route, user)
			req.SetPathValue(idPath, tt.id)

			if tt.setup != nil {
				req = tt.setup(t)
			}

			rr := httptest.NewRecorder()
			handler.DeleteByID(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(t, rs)
			}
		})
	}
}

func TestCategoryHandler_UpdateByID(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, svc, cleanup := setupTestCategoryHandler(t)
	defer cleanup()

	user := testutils.CreateTestUser(t, svc.User, "testuser")
	category, err := svc.Category.Create(&store.CreateCategoryParams{
		UserID: user.ID,
		Name:   "Test",
		Color:  "#123",
		Icon:   "T",
	})
	if err != nil {
		t.Fatal(err)
	}
	categoryID := strconv.Itoa(int(category.ID))

	route := "/v1/categories"
	idPath := "categoryID"

	tests := []struct {
		name           string
		id             string
		requestBody    any
		expectedStatus int
		setup          func(*testing.T) *http.Request
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name: "Update category",
			id:   categoryID,
			requestBody: map[string]any{
				"name":  "Test Update",
				"color": "#ABC",
				"icon":  "U",
			},
			checkResponse: func(t *testing.T, rs *http.Response) {
				var resBody map[string]*store.Category
				json.NewDecoder(rs.Body).Decode(&resBody)

				category := resBody["category"]
				assert.Equal(t, category.Name, "Test Update")
				assert.Equal(t, category.Color, "#ABC")
				assert.Equal(t, category.Icon, "U")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Update partial category",
			id:   categoryID,
			requestBody: map[string]any{
				"color": "#123456",
			},
			checkResponse: func(t *testing.T, rs *http.Response) {
				var resBody map[string]*store.Category
				json.NewDecoder(rs.Body).Decode(&resBody)

				category := resBody["category"]
				assert.Equal(t, category.Name, "Test Update")
				assert.Equal(t, category.Color, "#123456")
				assert.Equal(t, category.Icon, "U")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Fail to update other user's category",
			setup: func(t *testing.T) *http.Request {
				requestBody := map[string]any{
					"color": "#123456",
				}

				user2 := testutils.CreateTestUser(t, svc.User, "testuser2")
				category2 := testutils.CreateTestCategory(t, svc.Category, user2.ID)

				req := testutils.CreatePostRequest(t, route, requestBody, user)
				req.SetPathValue(idPath, strconv.Itoa(int(category2.ID)))

				return req
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Fail to update initial category",
			id:   "1",
			requestBody: map[string]any{
				"color": "#123456",
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Validation error",
			id:   categoryID,
			requestBody: map[string]any{
				"color": "FFF",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Incorrect JSON",
			id:             categoryID,
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Not found",
			id:             "999",
			requestBody:    map[string]any{},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid ID",
			id:             "bad",
			requestBody:    map[string]any{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty ID",
			id:             "",
			requestBody:    map[string]any{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Negative ID",
			id:             "-1",
			requestBody:    map[string]any{},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := testutils.CreatePostRequest(t, route, tt.requestBody, user)
			req.SetPathValue(idPath, tt.id)

			if tt.setup != nil {
				req = tt.setup(t)
			}

			rr := httptest.NewRecorder()
			handler.UpdateByID(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(t, rs)
			}
		})
	}
}
