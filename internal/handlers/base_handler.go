package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Quak1/gokei/internal/errors"
	"github.com/Quak1/gokei/internal/services"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type BaseHandler struct{}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

func (h *BaseHandler) GetJWTClaims(w http.ResponseWriter, r *http.Request) (*services.JWTCustomClaims, error) {
	claims, ok := r.Context().Value("claims").(*services.JWTCustomClaims)
	if !ok {
		return nil, errors.NewAppError(errors.ErrInternal, "", nil)
	}

	return claims, nil
}

func (h *BaseHandler) ParseRequest(w http.ResponseWriter, r *http.Request, dst any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.NewAppError(errors.ErrInternal, "Failed to read request data", err)
	}

	if len(body) == 0 {
		return errors.NewAppError(errors.ErrInternal, "Request body is required", err)
	}

	if err := json.Unmarshal(body, dst); err != nil {
		return errors.NewAppError(errors.ErrInternal, "Request body must be valid JSON", err)
	}

	if err := validate.Struct(dst); err != nil {
		return handleValidationErrors(err)
	}

	return nil
}

func handleValidationErrors(err error) error {
	fmt.Println("VALIDATION ERRORR TODO:", err)
	return errors.NewAppError(errors.ErrValidation, "Validation error ", err)
}
