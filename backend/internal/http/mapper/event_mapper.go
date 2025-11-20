package mapper

import (
	"time"

	"backend/internal/entities"
)

// Event DTOs

type EventResponse struct {
	ID                string    `json:"id"`
	Identifier        string    `json:"identifier"`
	Name              string    `json:"name"`
	Description       string    `json:"description,omitempty"`
	StartsAt          time.Time `json:"starts_at"`
	EndsAt            time.Time `json:"ends_at"`
	Timezone          string    `json:"timezone,omitempty"`
	LocationName      string    `json:"location_name,omitempty"`
	LogoURL           string    `json:"logo_url,omitempty"`
	ThumbnailImageURL string    `json:"thumbnail_image_url,omitempty"`
	LargeImageURL     string    `json:"large_image_url,omitempty"`
	OriginalImageURL  string    `json:"original_image_url,omitempty"`
	IconImageURL      string    `json:"icon_image_url,omitempty"`
	Privacy           string    `json:"privacy,omitempty"`
	State             string    `json:"state,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	Link              string    `json:"link,omitempty"`
}

// Mapping functions

func EventToResponse(event entities.Event) EventResponse {
	return EventResponse{
		ID:                event.ID,
		Identifier:        event.Identifier,
		Name:              event.Name,
		Description:       event.Description,
		StartsAt:          event.StartsAt,
		EndsAt:            event.EndsAt,
		Timezone:          event.Timezone,
		LocationName:      event.LocationName,
		LogoURL:           event.LogoURL,
		ThumbnailImageURL: event.ThumbnailImageURL,
		LargeImageURL:     event.LargeImageURL,
		OriginalImageURL:  event.OriginalImageURL,
		IconImageURL:      event.IconImageURL,
		Privacy:           event.Privacy,
		State:             event.State,
		CreatedAt:         event.CreatedAt,
		Link:              event.Link,
	}
}

func EventsToResponse(events []entities.Event) []EventResponse {
	result := make([]EventResponse, len(events))
	for i, event := range events {
		result[i] = EventToResponse(event)
	}
	return result
}
