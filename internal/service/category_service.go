package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/pkg/validator"
)

type CategoryService struct {
	queries store.QuerierTx
}

func NewCategoryService(queries store.QuerierTx) *CategoryService {
	return &CategoryService{
		queries: queries,
	}
}

func validateCategory(v *validator.Validator, category *store.Category) {
	v.Check(validator.NonZero(category.Name), "name", "Must be provided")
	v.Check(validator.MaxLength(category.Name, 20), "name", "Must not be more than 20 bytes long")

	v.Check(validator.NonZero(category.Color), "color", "Must be provided")
	v.Check(validator.HexColor(category.Color), "color", "Must be valid Hex Color")

	v.Check(validator.NonZero(category.Icon), "icon", "Must be provided")

	v.Check(validator.NonZero(category.UserID), "user_id", "Must be provided")
}

func (s *CategoryService) Create(categoryParams *store.CreateCategoryParams) (*store.Category, error) {
	category := &store.Category{
		UserID: categoryParams.UserID,
		Name:   categoryParams.Name,
		Color:  categoryParams.Color,
		Icon:   categoryParams.Icon,
	}

	v := validator.New()
	if validateCategory(v, category); !v.Valid() {
		return nil, v.GetErrors()
	}

	data, err := s.queries.CreateCategory(context.Background(), *categoryParams)
	if err != nil {
		return nil, database.HandleForeignKeyError(err)
	}

	return &data, nil
}

func (s *CategoryService) GetAll(userID int32) ([]*store.Category, error) {
	data, err := s.queries.GetAllCategories(context.Background(), store.GetAllCategoriesParams{
		AdminID: database.AdminUserID(),
		UserID:  userID,
	})
	if err != nil {
		return nil, err
	}

	categories := make([]*store.Category, len(data))
	for i, v := range data {
		categories[i] = &v
	}

	return categories, nil
}

func (s *CategoryService) GetByID(userID, categoryID int32) (*store.Category, error) {
	if userID < 1 || categoryID < 1 {
		return nil, database.ErrRecordNotFound
	}

	category, err := s.queries.GetCategoryByID(context.Background(), store.GetCategoryByIDParams{
		ID:     categoryID,
		UserID: userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (s *CategoryService) DeleteByID(userID, categoryID int32) error {
	if userID < 1 || categoryID < 1 {
		return database.ErrRecordNotFound
	}
	if categoryID == database.InitialCategoryID() {
		return database.ErrUpdateInitialCategory
	}

	result, err := s.queries.DeleteCategoryById(context.Background(), store.DeleteCategoryByIdParams{
		ID:     categoryID,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrRecordNotFound
	}

	return nil
}

type UpdateCategoryParams struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
	Icon  *string `json:"icon"`
}

func (s *CategoryService) UpdateByID(userId, categoryID int32, updateParams *UpdateCategoryParams) (*store.Category, error) {
	if userId < 1 || categoryID < 1 {
		return nil, database.ErrRecordNotFound
	}
	if categoryID == database.InitialCategoryID() {
		return nil, database.ErrUpdateInitialCategory
	}

	ctx := context.Background()

	category, err := s.queries.GetCategoryByID(ctx, store.GetCategoryByIDParams{
		ID:     categoryID,
		UserID: userId,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if updateParams.Name != nil {
		category.Name = *updateParams.Name
	}
	if updateParams.Color != nil {
		category.Color = *updateParams.Color
	}
	if updateParams.Icon != nil {
		category.Icon = *updateParams.Icon
	}

	v := validator.New()
	if validateCategory(v, &category); !v.Valid() {
		return nil, v.GetErrors()
	}

	result, err := s.queries.UpdateCategoryById(ctx, store.UpdateCategoryByIdParams{
		UserID:  category.UserID,
		Name:    category.Name,
		Color:   category.Color,
		Icon:    category.Icon,
		ID:      category.ID,
		Version: category.Version,
	})
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, database.ErrEditConflict
	}

	return &category, nil
}
