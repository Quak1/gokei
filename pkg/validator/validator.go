package validator

import (
	"regexp"
	"slices"
)

// from https://html.spec.whatwg.org/#valid-e-mail-address
var emailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var hexColorRX = regexp.MustCompile("^#([a-fA-F0-9]{6}|[a-fA-F0-9]{3})$")

type Validator struct {
	errors *ValidationError
}

func New() *Validator {
	return &Validator{
		errors: NewValidationError(),
	}
}

func (v *Validator) Valid() bool {
	return !v.errors.HasErrors()
}

func (v *Validator) AddError(key, message string) {
	v.errors.Add(key, message)
}

func (v *Validator) GetErrors() error {
	return v.errors
}

func (v *Validator) Check(ok bool, key, messate string) {
	if !ok {
		v.AddError(key, messate)
	}
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Email(email string) bool {
	return Matches(email, emailRX)
}

func HexColor(color string) bool {
	return Matches(color, hexColorRX)
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Unique[T comparable](values []T) bool {
	seen := make(map[T]struct{}, len(values))

	for _, value := range values {
		if _, exists := seen[value]; exists {
			return false
		}
		seen[value] = struct{}{}
	}

	return true
}

func NonZero[T comparable](value T) bool {
	var zero T
	return value != zero
}

func MaxLength(value string, size int) bool {
	return len(value) <= size
}

func MinLength(value string, size int) bool {
	return len(value) >= size
}
