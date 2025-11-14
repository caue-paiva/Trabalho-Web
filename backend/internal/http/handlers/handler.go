package handlers

import (
	"context"

	"backend/internal/platform/middleware"
	"backend/internal/server"

	"firebase.google.com/go/v4/auth"
)

type HandlerOption func(*BaseHandler)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	server     server.Server
	middleware []middleware.Middleware
}

// WithAuthMiddleware injects authorization middleware in the handler
func WithAuthMiddleware(ctx context.Context, authClient *auth.Client) HandlerOption {
	return func(handler *BaseHandler) {
		if handler != nil {
			handler.middleware = append(handler.middleware, middleware.NewAuthMiddleware(ctx, authClient))
		}
	}
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(srv server.Server, opts ...HandlerOption) *BaseHandler {
	handler := &BaseHandler{server: srv}
	for _, opt := range opts {
		opt(handler)
	}

	return handler
}
