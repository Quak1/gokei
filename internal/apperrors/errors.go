package apperrors

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	CodeInternal     ErrorCode = "INTERNAL"
	CodeNotFound     ErrorCode = "NOT_FOUND"
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	CodeForbidden    ErrorCode = "FORBIDDEN"
	CodeConflict     ErrorCode = "CONFLICT"
	CodeValidation   ErrorCode = "VALIDATION"
)

var codeToStatus = map[ErrorCode]int{
	CodeInternal:     http.StatusInternalServerError,
	CodeNotFound:     http.StatusNotFound,
	CodeUnauthorized: http.StatusUnauthorized,
	CodeForbidden:    http.StatusForbidden,
	CodeConflict:     http.StatusConflict,
	CodeValidation:   http.StatusUnprocessableEntity,
}

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) StatusCode() int {
	if status, ok := codeToStatus[e.Code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

func New(code ErrorCode, userMessage string, internalError error) *AppError {
	return &AppError{
		Code:    code,
		Message: userMessage,
		Err:     internalError,
	}
}

func NotFound(resource string, internalErr error) *AppError {
	return New(CodeNotFound, fmt.Sprintf("The requested %s was not found", resource), internalErr)
}
