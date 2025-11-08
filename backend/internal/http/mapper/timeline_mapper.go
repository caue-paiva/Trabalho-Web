package mapper

import (
	"fmt"
	"time"

	"backend/internal/entities"
)

// TimelineEntry DTOs

type CreateTimelineEntryRequest struct {
	Name     string `json:"name"`
	Text     string `json:"text"`
	Location string `json:"location,omitempty"`
	Date     string `json:"date"` // ISO format
}

type UpdateTimelineEntryRequest struct {
	Name     string `json:"name,omitempty"`
	Text     string `json:"text,omitempty"`
	Location string `json:"location,omitempty"`
	Date     string `json:"date,omitempty"` // ISO format
}

type TimelineEntryResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Text          string    `json:"text"`
	Location      string    `json:"location,omitempty"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastUpdatedBy string    `json:"last_updated_by,omitempty"`
}

// Mapping functions

func ToTimelineEntryEntity(req CreateTimelineEntryRequest) (entities.TimelineEntry, error) {
	// Parse date
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return entities.TimelineEntry{}, fmt.Errorf("invalid date format: %w", err)
	}

	return entities.TimelineEntry{
		Name:     req.Name,
		Text:     req.Text,
		Location: req.Location,
		Date:     date,
	}, nil
}

func ToTimelineEntryUpdateEntity(req UpdateTimelineEntryRequest) (entities.TimelineEntry, error) {
	var date time.Time
	var err error
	if req.Date != "" {
		date, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			return entities.TimelineEntry{}, fmt.Errorf("invalid date format: %w", err)
		}
	}

	return entities.TimelineEntry{
		Name:     req.Name,
		Text:     req.Text,
		Location: req.Location,
		Date:     date,
	}, nil
}

func TimelineEntryToResponse(entry entities.TimelineEntry) TimelineEntryResponse {
	return TimelineEntryResponse{
		ID:            entry.ID,
		Name:          entry.Name,
		Text:          entry.Text,
		Location:      entry.Location,
		Date:          entry.Date,
		CreatedAt:     entry.CreatedAt,
		UpdatedAt:     entry.UpdatedAt,
		LastUpdatedBy: entry.LastUpdatedBy,
	}
}

func TimelineEntriesToResponse(entries []entities.TimelineEntry) []TimelineEntryResponse {
	result := make([]TimelineEntryResponse, len(entries))
	for i, entry := range entries {
		result[i] = TimelineEntryToResponse(entry)
	}
	return result
}
