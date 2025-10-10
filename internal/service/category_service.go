package service

import (
	"context"

	"github.com/Quak1/gokei/internal/database/queries"
	"github.com/Quak1/gokei/pkg/validator"
)

type CategoryService struct {
	db *queries.Queries
}

func NewCategoryService(db *queries.Queries) *CategoryService {
	return &CategoryService{
		db: db,
	}
}

func validateCategory(v *validator.Validator, category *queries.CreateCategoryParams) {
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

func toPublicCategory(c queries.Category) *PublicCategory {
	return &PublicCategory{
		ID:    int(c.ID),
		Name:  c.Name,
		Color: c.Color,
		Icon:  c.Icon,
	}
}

func (s *CategoryService) Create(category *queries.CreateCategoryParams) (*PublicCategory, error) {
	v := validator.New()
	if validateCategory(v, category); !v.Valid() {
		return nil, v.GetErrors()
	}

	data, err := s.db.CreateCategory(context.Background(), *category)
	if err != nil {
		return nil, err
	}

	newCategory := toPublicCategory(data)

	return newCategory, nil
}

func (s *CategoryService) GetAll() ([]*PublicCategory, error) {
	data, err := s.db.GetAllCategories(context.Background())
	if err != nil {
		return nil, err
	}

	categories := make([]*PublicCategory, len(data))
	for i, v := range data {
		categories[i] = toPublicCategory(v)
	}

	return categories, nil
}
