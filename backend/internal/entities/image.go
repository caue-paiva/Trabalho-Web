package entities

import "time"

// Image represents image metadata
type Image struct {
	ID            string
	Slug          string    // Optional
	ObjectURL     string    // Storage URL
	Name          string
	Text          string    // Description
	Date          time.Time
	Location      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastUpdatedBy string
}
