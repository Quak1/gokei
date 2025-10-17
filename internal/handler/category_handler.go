package handler

import (
	"errors"
	"fmt"
	"net/http"

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
	var input store.CreateCategoryParams

	err := response.ReadJSON(w, r, &input)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	category, err := h.categoryService.Create(&input)
	if err != nil {
		var validationErr *validator.ValidationError

		switch {
		case errors.As(err, &validationErr):
			response.FailedValidationResponse(w, r, validationErr)
		default:
			// TODO handle errors
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
	categories, err := h.categoryService.GetAll()

	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}

	err = response.OK(w, response.Envelope{"categories": categories})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "categoryID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	category, err := h.categoryService.GetByID(int32(id))
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
	id, err := readIntParam(r, "categoryID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	err = h.categoryService.DeleteByID(int32(id))
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
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
