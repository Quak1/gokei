package apperrors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	*AppError
	Fields map[string]string
}

func NewValidation(message string, err error) *ValidationError {
	return &ValidationError{
		AppError: New(CodeValidation, message, nil),
		Fields:   handleValidationFields(err),
	}
}

func handleValidationFields(err error) map[string]string {
	fields := map[string]string{}

	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		for _, e := range validateErrs {
			field := strings.ToLower(e.Field())
			fields[field] = getValidationMessage(e)
		}
	}

	return fields
}

func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "Field is required"
	case "oneof":
		return fmt.Sprintf("Field must be one of: %s", e.Param())
	default:
		return "Field is invalid"
	}
}
