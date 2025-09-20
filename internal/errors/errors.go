package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrInternal     ErrorCode = "INTERNAL"
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrForbidden    ErrorCode = "FORBIDDEN"
	ErrConflict     ErrorCode = "CONFLICT"
	ErrValidation   ErrorCode = "VALIDATION"
)

var errorToStatus = map[ErrorCode]int{
	ErrInternal:     http.StatusInternalServerError,
	ErrNotFound:     http.StatusNotFound,
	ErrUnauthorized: http.StatusUnauthorized,
	ErrForbidden:    http.StatusForbidden,
	ErrConflict:     http.StatusConflict,
	ErrValidation:   http.StatusUnprocessableEntity,
}

type AppError struct {
	Code        ErrorCode
	UserMessage string
	InternalErr error
	StatusCode  int
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Code, e.UserMessage, e.InternalErr)
}

func (e *AppError) Unwrap() error {
	return e.InternalErr
}

func NewAppError(code ErrorCode, userMessage string, internalError error) *AppError {
	return &AppError{
		Code:        code,
		UserMessage: userMessage,
		InternalErr: internalError,
		StatusCode:  errorToStatus[code],
	}
}

func NewNotFoundError(resource string, internalErr error) *AppError {
	return NewAppError(ErrNotFound, fmt.Sprintf("The requested %s was not found", resource), internalErr)
}
