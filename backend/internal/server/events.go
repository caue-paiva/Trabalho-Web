package server

import (
	"context"
	"fmt"

	"backend/internal/entities"
)

const (
	grupyBaseEventsWebPageURL = "https://eventos.grupysanca.com.br"
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
	events, err := s.events.GetEvents(ctx, limit, orderBy, desc)
	if err != nil {
		return nil, err
	}
	addLinksToevents(events)
	return events, nil
}

// fills the Link field of events in place
func addLinksToevents(events []entities.Event) {
	for i := range events {
		events[i].Link = fmt.Sprintf("%s/e/%s", grupyBaseEventsWebPageURL, events[i].ID)
	}
}
