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
	v.Check(validator.MaxLength(category.Name, 20), "name", "Must not be more than 10 bytes long")

	v.Check(validator.NonZero(category.Color), "color", "Must be provided")
	v.Check(validator.HexColor(category.Color), "color", "Must be valid Hex Color")

	v.Check(validator.NonZero(category.Icon), "icon", "Must be provided")
}

type PublicCategory struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

func toPublicCategory(c store.Category) *PublicCategory {
	return &PublicCategory{
		ID:    int(c.ID),
		Name:  c.Name,
		Color: c.Color,
		Icon:  c.Icon,
	}
}

func (s *CategoryService) Create(category *store.CreateCategoryParams) (*PublicCategory, error) {
	v := validator.New()
	if validateCategory(v, category); !v.Valid() {
		return nil, v.GetErrors()
	}

	data, err := s.queries.CreateCategory(context.Background(), *category)
	if err != nil {
		return nil, err
	}

	newCategory := toPublicCategory(data)

	return newCategory, nil
}

func (s *CategoryService) GetAll() ([]*PublicCategory, error) {
	data, err := s.queries.GetAllCategories(context.Background())
	if err != nil {
		return nil, err
	}

	categories := make([]*PublicCategory, len(data))
	for i, v := range data {
		categories[i] = toPublicCategory(v)
	}

	return categories, nil
}
