package service

import (
	"context"

	"backend/internal/entities"
)

// =======================
// EVENTS OPERATIONS
// =======================

func (s *server) GetEvents(ctx context.Context, limit int, orderBy string, desc bool) ([]entities.Event, error) {
	// Business logic: validate and set defaults
	if limit <= 0 || limit > 100 {
		limit = 10 // default
	}

	if orderBy == "" {
		orderBy = "startDate"
	}

	// Delegate to port
	return s.events.GetEvents(ctx, limit, orderBy, desc)
}
