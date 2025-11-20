package mapper

import (
	"time"

	"backend/internal/entities"
)

// GaleryEvent DTOs

// CreateGaleryEventRequest represents the request to create a galery event
type CreateGaleryEventRequest struct {
	Name         string    `json:"name" binding:"required"`
	Location     string    `json:"location" binding:"required"`
	Date         time.Time `json:"date" binding:"required"`
	ImagesBase64 []string  `json:"images_base64" binding:"required,min=1"`
}

// GaleryEventResponse represents a galery event response
type GaleryEventResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Date      time.Time `json:"date"`
	ImageURLs []string  `json:"image_urls"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Mapping functions

// GaleryEventToResponse converts a GaleryEvent entity to a response DTO
func GaleryEventToResponse(event entities.GaleryEvent) GaleryEventResponse {
	return GaleryEventResponse{
		ID:        event.ID,
		Name:      event.Name,
		Location:  event.Location,
		Date:      event.Date,
		ImageURLs: event.ImageURLs,
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}
}

// GaleryEventsToResponse converts multiple GaleryEvent entities to response DTOs
func GaleryEventsToResponse(events []entities.GaleryEvent) []GaleryEventResponse {
	result := make([]GaleryEventResponse, len(events))
	for i, event := range events {
		result[i] = GaleryEventToResponse(event)
	}
	return result
}
