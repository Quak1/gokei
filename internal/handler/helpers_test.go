package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Quak1/gokei/pkg/assert"
)

func Test_ReadIntParam(t *testing.T) {
	tests := []struct {
		name         string
		setupRequest func(*testing.T, string) *http.Request
		key          string
		wantError    bool
		expected     int
	}{
		{
			name: "Get number",
			setupRequest: func(t *testing.T, key string) *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.SetPathValue(key, strconv.Itoa(10))
				return r
			},
			key:       "num",
			wantError: false,
			expected:  10,
		},
		{
			name: "Missing key",
			setupRequest: func(t *testing.T, key string) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			key:       "num",
			wantError: true,
		},
		{
			name: "Value is not a number",
			setupRequest: func(t *testing.T, key string) *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.SetPathValue(key, "NaN")
				return r
			},
			key:       "num",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setupRequest(t, tt.key)
			num, err := readIntParam(r, tt.key)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Equal(t, num, tt.expected)
		})
	}
}
