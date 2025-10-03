package entities

import "time"

// TimelineEntry represents a timeline event
type TimelineEntry struct {
	ID            string
	Name          string
	Text          string
	Location      string
	Date          time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastUpdatedBy string
}
