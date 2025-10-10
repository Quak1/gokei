package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Quak1/gokei/internal/database/queries"
	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/pkg/response"
	"github.com/Quak1/gokei/sql/validator"
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
	var input queries.CreateCategoryParams

	err := response.ReadJSON(w, r, &input)
	if err != nil {
		response.BadRequestResponse(w, r)
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
		response.BadRequestResponse(w, r)
	}
}
