package validator

import (
	"testing"

	"github.com/Quak1/gokei/pkg/assert"
)

func TestCheck(t *testing.T) {
	v := New()

	v.Check(true, "field1", "should not add error")
	assert.Equal(t, v.Valid(), true)

	v.Check(false, "field2", "should add error")
	v.Check(false, "field3", "another error")
	assert.Equal(t, v.Valid(), false)

	errors := v.GetErrors().(*ValidationError)
	assert.Equal(t, len(errors.Errors), 2)
	assert.Equal(t, errors.Errors["field2"], "should add error")
	assert.Equal(t, errors.Errors["field3"], "another error")
}

func TestEmail(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid email", "test@test.com", true},
		{"invalid", "test", false},
		{"invalid - empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Email(tt.input)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestHexColor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid 6-digit hex", "#ABC123", true},
		{"valid 3-digit hex", "#DeF", true},
		{"invalid - no hash", "FF5733", false},
		{"invalid - wrong chars", "#GGG", false},
		{"invalid - wrong length", "#FF", false},
		{"invalid - empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HexColor(tt.input)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestPermittedValue(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		values []string
		want   bool
	}{
		{"valid", "test", []string{"test", "test1"}, true},
		{"value not valid", "test", []string{"test1", "test2"}, false},
		{"no valuse", "test", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PermittedValue(tt.input, tt.values...)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  bool
	}{
		{"all unique values", []string{"test", "test2"}, true},
		{"no values", []string{}, true},
		{"duplicate value", []string{"test", "test2", "test"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.input)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestNonZero(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"non zero string", "test", true},
		{"zero string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NonZero(tt.input)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int
		want   bool
	}{
		{"valid", "test", 4, true},
		{"valid empty input", "", 0, true},
		{"over max length", "test", 2, false},
		{"negative max length", "", -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxLength(tt.input, tt.length)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int
		want   bool
	}{
		{"valid", "test", 4, true},
		{"valid empty input", "", 0, true},
		{"under min length", "test", 5, false},
		{"negative min length", "", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinLength(tt.input, tt.length)
			assert.Equal(t, got, tt.want)
		})
	}
}
