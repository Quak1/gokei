package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/Quak1/gokei/internal/apperrors"
	"github.com/Quak1/gokei/internal/services"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
}

type BaseHandler struct{}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

func (h *BaseHandler) GetJWTClaims(w http.ResponseWriter, r *http.Request) (*services.JWTCustomClaims, error) {
	claims, ok := r.Context().Value("claims").(*services.JWTCustomClaims)
	if !ok {
		return nil, apperrors.New(apperrors.CodeInternal, "", nil)
	}

	return claims, nil
}

func (h *BaseHandler) ParseRequest(w http.ResponseWriter, r *http.Request, dst any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return apperrors.New(apperrors.CodeInternal, "Failed to read request data", err)
	}

	if len(body) == 0 {
		return apperrors.New(apperrors.CodeInternal, "Request body is required", err)
	}

	if err := json.Unmarshal(body, dst); err != nil {
		return apperrors.New(apperrors.CodeInternal, "Request body must be valid JSON", err)
	}

	if err := validate.Struct(dst); err != nil {
		return apperrors.NewValidation("Validation error", err)
	}

	return nil
}
