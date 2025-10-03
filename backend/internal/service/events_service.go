package service

import (
	"context"

	"backend/internal/entities"
)

// EventsService defines business logic for events operations
type EventsService interface {
	GetEvents(ctx context.Context, limit int, orderBy string, desc bool) ([]entities.Event, error)
}

// eventsService implements EventsService
type eventsService struct {
	events GrupyEventsPort
}

// NewEventsService creates a new EventsService
func NewEventsService(events GrupyEventsPort) EventsService {
	return &eventsService{
		events: events,
	}
}

func (s *eventsService) GetEvents(ctx context.Context, limit int, orderBy string, desc bool) ([]entities.Event, error) {
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
