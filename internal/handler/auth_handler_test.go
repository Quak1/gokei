package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/internal/testutils"
	"github.com/Quak1/gokei/pkg/assert"
)

func TestAuthHandler_Login(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, cleanup, err := testutils.NewTestDB()
	if err != nil {
		cleanup()
		t.Fatal(err)
	}
	defer cleanup()

	svc := service.New(db)
	handler := NewAuthHandler(svc.Auth)
	user := testutils.CreateTestUser(t, svc.User, "testuser")
	route := "/v1/auth/login"

	tests := []struct {
		name           string
		input          any
		expectedStatus int
		validate       func(*testing.T, *http.Response)
	}{
		{
			name: "Login",
			input: map[string]any{
				"username": "testuser",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]*service.Token
				json.NewDecoder(r.Body).Decode(&resBody)
				token := resBody["authentication_token"]

				tokenUser, err := svc.User.GetForToken(token.Plaintext)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tokenUser.ID, user.ID)
			},
		},
		{
			name: "Wrong password",
			input: map[string]any{
				"username": "testuser",
				"password": "password",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Not found user",
			input: map[string]any{
				"username": "user",
				"password": "password",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Validation error - short username",
			input: map[string]any{
				"username": "u",
				"password": "password",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]map[string]any
				json.NewDecoder(r.Body).Decode(&resBody)

				errors := resBody["error"]
				assert.Equal(t, len(errors), 1)
				if _, ok := errors["username"]; !ok {
					t.Error("Missing username validation error message")
				}
			},
		},
		{
			name: "Validation error - short password",
			input: map[string]any{
				"username": "user",
				"password": "pass",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			validate: func(t *testing.T, r *http.Response) {
				var resBody map[string]map[string]any
				json.NewDecoder(r.Body).Decode(&resBody)

				errors := resBody["error"]
				assert.Equal(t, len(errors), 1)
				if _, ok := errors["password"]; !ok {
					t.Error("Missing password validation error message")
				}
			},
		},
		{
			name:           "Empty JSON",
			input:          "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
			req.Header.Set("Contenty-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.Login(rr, req)

			rs := rr.Result()
			defer rs.Body.Close()

			assert.Equal(t, rs.StatusCode, tt.expectedStatus)

			if tt.validate != nil {
				tt.validate(t, rs)
			}
		})
	}
}
