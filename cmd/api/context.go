package main

import (
	"context"
	"net/http"

	"github.com/namikaze-dev/bluelight/internal/models"
)

type contextKey string

const userContextKey = contextKey("user")

// contextSetUse returns a new copy of the request with the provided
// User struct added to the context.
func (app *application) contextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser retrieves the User struct from the request context.
func (app *application) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
