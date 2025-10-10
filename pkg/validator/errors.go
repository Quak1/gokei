package validator

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Errors map[string]string
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors: make(map[string]string),
	}
}

func (v ValidationError) Error() string {
	if len(v.Errors) == 0 {
		return "validation failed"
	}

	var messages []string
	for field, msg := range v.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", field, msg))
	}

	return "validation failed: " + strings.Join(messages, "; ")
}

func (v ValidationError) HasErrors() bool {
	return len(v.Errors) > 0
}

func (v *ValidationError) Add(field, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}
	v.Errors[field] = message
}
