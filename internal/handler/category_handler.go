package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Quak1/gokei/internal/appcontext"
	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/database/store"
	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/pkg/response"
	"github.com/Quak1/gokei/pkg/validator"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: svc,
	}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Color string `json:"color"`
		Icon  string `json:"icon"`
	}

	err := response.ReadJSON(w, r, &input)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	params := store.CreateCategoryParams{
		UserID: ctxUser.ID,
		Name:   input.Name,
		Color:  input.Color,
		Icon:   input.Icon,
	}

	category, err := h.categoryService.Create(&params)
	if err != nil {
		var validationErr *validator.ValidationError

		switch {
		case errors.As(err, &validationErr):
			response.FailedValidationResponse(w, r, validationErr)
		case errors.Is(err, database.ErrInvalidUser):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}

		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/categories/%d", category.ID))

	err = response.Created(w, response.Envelope{"category": category}, headers)
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctxUser := appcontext.GetContextUser(r)

	categories, err := h.categoryService.GetAll(ctxUser.ID)
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}

	err = response.OK(w, response.Envelope{"categories": categories})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := readIntParam(r, "categoryID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	category, err := h.categoryService.GetByID(ctxUser.ID, int32(categoryID))
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"category": category})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *CategoryHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := readIntParam(r, "categoryID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	err = h.categoryService.DeleteByID(ctxUser.ID, int32(categoryID))
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		case errors.Is(err, database.ErrUpdateInitialCategory):
			response.ForbiddenResponse(w, r, err)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"message": "category successfully deleted"})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *CategoryHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := readIntParam(r, "categoryID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	var input service.UpdateCategoryParams
	err = response.ReadJSON(w, r, &input)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	category, err := h.categoryService.UpdateByID(ctxUser.ID, int32(categoryID), &input)
	if err != nil {
		var validationErr *validator.ValidationError
		switch {
		case errors.As(err, &validationErr):
			response.FailedValidationResponse(w, r, validationErr)
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		case errors.Is(err, database.ErrEditConflict):
			response.ConflictResponse(w, r)
		case errors.Is(err, database.ErrUpdateInitialCategory):
			response.ForbiddenResponse(w, r, err)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"category": category})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}
