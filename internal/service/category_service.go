package service

import (
	"context"

	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/pkg/validator"
)

type CategoryService struct {
	queries *store.Queries
}

func NewCategoryService(queries *store.Queries) *CategoryService {
	return &CategoryService{
		queries: queries,
	}
}

func validateCategory(v *validator.Validator, category *store.CreateCategoryParams) {
	v.Check(validator.NonZero(category.Name), "name", "Must be provided")
	v.Check(validator.MaxLength(category.Name, 20), "name", "Must not be more than 20 bytes long")

	v.Check(validator.NonZero(category.Color), "color", "Must be provided")
	v.Check(validator.HexColor(category.Color), "color", "Must be valid Hex Color")

	v.Check(validator.NonZero(category.Icon), "icon", "Must be provided")
}

func (s *CategoryService) Create(category *store.CreateCategoryParams) (*store.Category, error) {
	v := validator.New()
	if validateCategory(v, category); !v.Valid() {
		return nil, v.GetErrors()
	}

	data, err := s.queries.CreateCategory(context.Background(), *category)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *CategoryService) GetAll() ([]*store.Category, error) {
	data, err := s.queries.GetAllCategories(context.Background())
	if err != nil {
		return nil, err
	}

	categories := make([]*store.Category, len(data))
	for i, v := range data {
		categories[i] = &v
	}

	return categories, nil
}
