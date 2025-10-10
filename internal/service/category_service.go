package service

import (
	"context"

	"github.com/Quak1/gokei/internal/database/queries"
	"github.com/Quak1/gokei/sql/validator"
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

func (s *CategoryService) Create(category *queries.CreateCategoryParams) (*queries.CreateCategoryRow, error) {
	v := validator.New()
	if validateCategory(v, category); !v.Valid() {
		return nil, v.GetErrors()
	}

	data, err := s.db.CreateCategory(context.Background(), *category)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
