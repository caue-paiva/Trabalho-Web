package server

import (
	"context"

	"backend/internal/entities"
)

// =======================
// EVENTS OPERATIONS
// =======================

func (s *server) GetEvents(ctx context.Context, limit int, orderBy string, desc bool) ([]entities.Event, error) {
	// Validate limit
	if limit <= 0 || limit > 100 {
		limit = 10 // default
	}

	// Delegate to port 
	return s.events.GetEvents(ctx, limit, orderBy, desc)
}
