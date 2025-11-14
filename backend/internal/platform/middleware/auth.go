package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	authcfg "backend/internal/platform/auth"
	customerrors "backend/internal/platform/errors"

	"firebase.google.com/go/v4/auth"
)

// Middleware is a func which takes in an http request and returns it and an error
type Middleware func(*http.Request) (*http.Request, error)

func NewAuthMiddleware(ctx context.Context, authClient *auth.Client) Middleware {
	return func(r *http.Request) (*http.Request, error) {
		if r == nil {
			return nil, errors.New("nil http request")
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return r, customerrors.ErrUnauthorized
		}

		idToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if idToken == "" {
			return r, customerrors.ErrUnauthorized
		}

		// we will not use the idToken
		_, err := authClient.VerifyIDToken(ctx, idToken)
		if err != nil {
			return r, customerrors.ErrUnauthorized
		}

		return r, nil
	}
}

func NewAuthMiddlewareFunc(nextHandle func(w http.ResponseWriter, r *http.Request), authCfg authcfg.AuthConfig) func(w http.ResponseWriter, r *http.Request) {
	if authCfg.Client == nil {
		return func(w http.ResponseWriter, r *http.Request) {
			nextHandle(w, r)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		idToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if idToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// we will not use the idToken
		_, err := authCfg.Client.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		nextHandle(w, r)
	}
}
