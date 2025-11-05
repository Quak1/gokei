package handler

import (
	"context"
	"net/http"

	"github.com/Quak1/gokei/internal/database/store"
)

type contextKey string

const userContextKey = contextKey("user")

func setContextUser(r *http.Request, user *store.GetUserFromTokenRow) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func getContextUser(r *http.Request) *store.GetUserFromTokenRow {
	user, ok := r.Context().Value(userContextKey).(*store.GetUserFromTokenRow)
	if !ok {
		panic("missing user in request context")
	}

	return user
}
