package appcontext

import (
	"context"
	"net/http"

	"github.com/Quak1/gokei/internal/database/store"
)

type contextKey string

const userContextKey = contextKey("user")

func SetContextUser(r *http.Request, user *store.GetUserFromTokenRow) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetContextUser(r *http.Request) *store.GetUserFromTokenRow {
	user, ok := r.Context().Value(userContextKey).(*store.GetUserFromTokenRow)
	if !ok {
		panic("missing user in request context")
	}

	return user
}
