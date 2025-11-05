package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Quak1/gokei/internal/appcontext"
	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/pkg/response"
	"github.com/Quak1/gokei/pkg/validator"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		userService: svc,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input service.InputUser

	err := response.ReadJSON(w, r, &input)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.userService.Create(&input)
	if err != nil {
		var validationErr *validator.ValidationError

		switch {
		case errors.As(err, &validationErr):
			response.FailedValidationResponse(w, r, validationErr)
		case errors.Is(err, service.ErrDuplicateUsername):
			response.BadRequestResponse(w, r, fmt.Errorf("Username is already in use"))
		default:
			response.ServerErrorResponse(w, r, err)
		}

		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/users/%d", user.ID))

	err = response.Created(w, response.Envelope{"user": user}, headers)
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "userID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	user, err := h.userService.GetByID(int32(id))
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"user": user})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *UserHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "userID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	ctxUser := appcontext.GetContextUser(r)

	if int(ctxUser.ID) != id {
		response.ForbiddenResponse(w, r)
		return
	}

	err = h.userService.DeleteByID(int32(id))
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"message": "user successfully deleted"})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *UserHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id, err := readIntParam(r, "userID")
	if err != nil {
		response.BadRequestResponseGeneric(w, r)
		return
	}

	var input service.UpdateUserParams
	err = response.ReadJSON(w, r, &input)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.userService.UpdateByID(int32(id), &input)
	if err != nil {
		var validationErr *validator.ValidationError
		switch {
		case errors.As(err, &validationErr):
			response.FailedValidationResponse(w, r, validationErr)
		case errors.Is(err, database.ErrRecordNotFound):
			response.NotFoundResponse(w, r)
		case errors.Is(err, database.ErrEditConflict):
			response.ConflictResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"user": user})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	token, err := h.userService.CreateAuthToken(input.Username, input.Password)
	if err != nil {
		var validationErr *validator.ValidationError
		switch {
		case errors.As(err, &validationErr):
			response.FailedValidationResponse(w, r, validationErr)
		case errors.Is(err, database.ErrRecordNotFound), errors.Is(err, service.ErrInvalidCredentials):
			response.InvalidCredentialsResponse(w, r)
		default:
			response.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = response.OK(w, response.Envelope{"authentication_token": token})
	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}
