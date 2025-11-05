package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Quak1/gokei/internal/appcontext"
	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/service"
	"github.com/Quak1/gokei/pkg/response"
	"github.com/Quak1/gokei/pkg/validator"
)

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			response.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			response.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()
		if service.ValidateTokenPlaintext(v, token); !v.Valid() {
			response.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := m.service.User.GetForToken(token)
		if err != nil {
			switch {
			case errors.Is(err, database.ErrRecordNotFound):
				response.InvalidAuthenticationTokenResponse(w, r)
			default:
				response.ServerErrorResponse(w, r, err)
			}
			return
		}

		r = appcontext.SetContextUser(r, user)

		next.ServeHTTP(w, r)
	})
}
