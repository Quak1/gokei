package handler

import (
	"bytes"
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

func setupTestCategoryHandler(t *testing.T) (*CategoryHandler, *service.CategoryService, func()) {
	db, cleanup, err := testutils.NewTestDB()
	if err != nil {
		t.Fatalf("test db setup failed: %v", err)
	}

	svc := service.New(db)
	handler := NewCategoryHandler(svc.Category)

	return handler, svc.Category, cleanup
}

func TestCategoryHandler_Create(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	handler, _, cleanup := setupTestCategoryHandler(t)
	defer cleanup()

	route := "/v1/categories"

	tests := []struct {
		name           string
		requestBody    any
		expectedStatus int
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
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
			req.Header.Set("Contenty-Type", "application/json")

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

	handler, _, cleanup := setupTestCategoryHandler(t)
	defer cleanup()

	route := "/v1/categories"

	tests := []struct {
		name           string
		expectedStatus int
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, route, nil)

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

	handler, _, cleanup := setupTestCategoryHandler(t)
	defer cleanup()

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
			id:             "1",
			expectedStatus: http.StatusOK,
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
			req := httptest.NewRequest(http.MethodGet, route, nil)
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

	category, err := svc.Create(&store.CreateCategoryParams{Name: "Test", Color: "#123", Icon: "T"})
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
		checkResponse  func(*testing.T, *http.Response)
	}{
		{
			name:           "Delete category",
			id:             categoryID,
			expectedStatus: http.StatusOK,
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
			req := httptest.NewRequest(http.MethodGet, route, nil)
			req.SetPathValue(idPath, tt.id)

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

	category, err := svc.Create(&store.CreateCategoryParams{Name: "Test", Color: "#123", Icon: "T"})
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
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
			req.Header.Set("Contenty-Type", "application/json")
			req.SetPathValue(idPath, tt.id)

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
