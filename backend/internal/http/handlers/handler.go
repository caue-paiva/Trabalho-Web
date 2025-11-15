package handlers

import (
	"backend/internal/platform/middleware"
	"backend/internal/server"
)

type HandlerOption func(*BaseHandler)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	server     server.Server
	middleware []middleware.Middleware
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(srv server.Server, opts ...HandlerOption) *BaseHandler {
	handler := &BaseHandler{server: srv}
	for _, opt := range opts {
		opt(handler)
	}

	return handler
}
