package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/pkg/assert"
)

func newTestService(mock func(m *store.MockQuerierTx)) *CategoryService {
	queriesMock := &store.MockQuerierTx{}
	mock(queriesMock)
	return NewCategoryService(queriesMock)
}

func checkError(t *testing.T, wantErr bool, err error) {
	if wantErr {
		if err == nil {
			t.Error("Expected error, got nil")
		}
	} else {
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}

func Test_Create(t *testing.T) {
	validCategory := store.CreateCategoryParams{
		Name:  "Travel",
		Color: "#FFF",
		Icon:  "T",
	}
	tests := []struct {
		name    string
		input   store.CreateCategoryParams
		mock    func(*store.MockQuerierTx)
		wantErr bool
	}{
		{
			name:    "Successful creation",
			input:   validCategory,
			mock:    func(m *store.MockQuerierTx) {},
			wantErr: false,
		},
		{
			name:  "Database error",
			input: validCategory,
			mock: func(m *store.MockQuerierTx) {
				m.CreateCategoryFunc = func(ctx context.Context, arg store.CreateCategoryParams) (store.Category, error) {
					return store.Category{}, errors.New("database error")
				}
			},
			wantErr: true,
		},
		{
			name: "Validation error - missing name",
			input: store.CreateCategoryParams{
				Name:  "",
				Color: "#FFF",
				Icon:  "T",
			},
			mock:    func(m *store.MockQuerierTx) {},
			wantErr: true,
		},
		{
			name: "Validation error - invalid color",
			input: store.CreateCategoryParams{
				Name:  "Travel",
				Color: "color",
				Icon:  "T",
			},
			mock:    func(m *store.MockQuerierTx) {},
			wantErr: true,
		},
		{
			name: "Validation error - missing icon",
			input: store.CreateCategoryParams{
				Name:  "Travel",
				Color: "#FFF",
				Icon:  "",
			},
			mock:    func(m *store.MockQuerierTx) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.mock)
			_, err := svc.Create(&tt.input)
			checkError(t, tt.wantErr, err)
		})
	}
}

func Test_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(*store.MockQuerierTx)
		wantErr bool
	}{
		{
			name:    "OK",
			mock:    func(m *store.MockQuerierTx) {},
			wantErr: false,
		},
		{
			name: "Database error",
			mock: func(m *store.MockQuerierTx) {
				m.GetAllCategoriesFunc = func(ctx context.Context) ([]store.Category, error) {
					return []store.Category{}, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.mock)
			_, err := svc.GetAll()
			checkError(t, tt.wantErr, err)
		})
	}
}

func Test_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		input   int32
		mock    func(*store.MockQuerierTx)
		wantErr bool
	}{
		{
			name:    "OK",
			input:   1,
			mock:    func(m *store.MockQuerierTx) {},
			wantErr: false,
		},
		{
			name:    "Invalid ID",
			input:   0,
			mock:    func(m *store.MockQuerierTx) {},
			wantErr: true,
		},
		{
			name:  "Database error",
			input: 1,
			mock: func(m *store.MockQuerierTx) {
				m.GetCategoryByIDFunc = func(ctx context.Context, id int32) (store.Category, error) {
					return store.Category{}, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.mock)
			_, err := svc.GetByID(tt.input)
			checkError(t, tt.wantErr, err)
		})
	}
}

func Test_DeleteByID(t *testing.T) {
	tests := []struct {
		name    string
		input   int32
		mock    func(*store.MockQuerierTx)
		wantErr bool
	}{
		{
			name:    "OK",
			input:   1,
			mock:    func(m *store.MockQuerierTx) {},
			wantErr: false,
		},
		{
			name:    "Invalid ID",
			input:   0,
			mock:    func(m *store.MockQuerierTx) {},
			wantErr: true,
		},
		{
			name:  "Database error",
			input: 1,
			mock: func(m *store.MockQuerierTx) {
				m.DeleteCategoryByIdFunc = func(ctx context.Context, id int32) (sql.Result, error) {
					return store.NewMockResult(0), errors.New("database error")
				}
			},
			wantErr: true,
		},
		{
			name:  "Not found ID",
			input: 1,
			mock: func(m *store.MockQuerierTx) {
				m.DeleteCategoryByIdFunc = func(ctx context.Context, id int32) (sql.Result, error) {
					return store.NewMockResult(0), errors.New("not found error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.mock)
			err := svc.DeleteByID(tt.input)
			checkError(t, tt.wantErr, err)
		})
	}
}

func getStringPointer(s string) *string {
	return &s
}

func Test_UpdateByID(t *testing.T) {
	validCategory := UpdateCategoryParams{
		Name:  getStringPointer("Test"),
		Color: getStringPointer("#FFF"),
		Icon:  getStringPointer("T"),
	}

	tests := []struct {
		name    string
		input   UpdateCategoryParams
		id      int32
		mock    func(*store.MockQuerierTx)
		expect  *store.Category
		wantErr bool
	}{
		{
			name:  "OK",
			input: validCategory,
			id:    1,
			mock: func(m *store.MockQuerierTx) {
				m.GetCategoryByIDFunc = func(ctx context.Context, id int32) (store.Category, error) {
					return store.Category{Name: "Original", Color: "#000", Icon: "O"}, nil
				}
			},
			expect:  &store.Category{Name: "Test", Color: "#FFF", Icon: "T"},
			wantErr: false,
		},
		{
			name:  "OK no update",
			input: UpdateCategoryParams{},
			id:    1,
			mock: func(m *store.MockQuerierTx) {
				m.GetCategoryByIDFunc = func(ctx context.Context, id int32) (store.Category, error) {
					return store.Category{Name: "Original", Color: "#000", Icon: "O"}, nil
				}
			},
			expect:  &store.Category{Name: "Original", Color: "#000", Icon: "O"},
			wantErr: false,
		},
		{
			name:  "validation error - invalid color",
			input: UpdateCategoryParams{Color: getStringPointer("red")},
			id:    1,
			mock: func(m *store.MockQuerierTx) {
				m.GetCategoryByIDFunc = func(ctx context.Context, id int32) (store.Category, error) {
					return store.Category{Name: "Original", Color: "#000", Icon: "O"}, nil
				}
			},
			wantErr: true,
		},
		{
			name:    "Invalid ID",
			input:   validCategory,
			id:      0,
			mock:    func(m *store.MockQuerierTx) {},
			wantErr: true,
		},
		{
			name:  "error id not found",
			input: validCategory,
			id:    1,
			mock: func(m *store.MockQuerierTx) {
				m.GetCategoryByIDFunc = func(ctx context.Context, id int32) (store.Category, error) {
					return store.Category{}, errors.New("not found error")
				}
			},
			wantErr: true,
		},
		{
			name:  "error no update done",
			input: validCategory,
			id:    1,
			mock: func(m *store.MockQuerierTx) {
				m.UpdateCategoryByIdFunc = func(ctx context.Context, arg store.UpdateCategoryByIdParams) (sql.Result, error) {
					return store.NewMockResult(0), nil
				}
			},
			wantErr: true,
		},
		{
			name:  "database error",
			input: validCategory,
			id:    0,
			mock: func(m *store.MockQuerierTx) {
				m.UpdateCategoryByIdFunc = func(ctx context.Context, arg store.UpdateCategoryByIdParams) (sql.Result, error) {
					return store.NewMockResult(1), errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.mock)
			data, err := svc.UpdateByID(tt.id, &tt.input)
			checkError(t, tt.wantErr, err)

			if tt.expect != nil {
				assert.Equal(t, data.Name, tt.expect.Name)
				assert.Equal(t, data.Color, tt.expect.Color)
				assert.Equal(t, data.Icon, tt.expect.Icon)
			}
		})
	}
}
