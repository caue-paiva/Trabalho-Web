package clients

import (
	"context"

	"backend/internal/entities"
	"backend/internal/gateway"
	"backend/internal/service"
)

// Compile-time interface check
var _ service.GrupyEventsPort = (*eventsClient)(nil)

type eventsClient struct {
	gateway *gateway.GrupyEventsAPI
}

// NewEventsClient creates a new GrupyEventsPort implementation
func NewEventsClient() service.GrupyEventsPort {
	return &eventsClient{
		gateway: gateway.NewGrupyEventsAPI(),
	}
}

func (c *eventsClient) GetEvents(ctx context.Context, limit int, orderBy string, desc bool) ([]entities.Event, error) {
	return c.gateway.GetEvents(ctx, limit, orderBy, desc)
}
