package handler

import (
	"errors"
	"net/http"

	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/pkg/response"
	"github.com/Quak1/gokei/pkg/validator"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: svc,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := response.ReadJSON(w, r, &input); err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	token, err := h.AuthService.CreateAuthToken(input.Username, input.Password)
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
