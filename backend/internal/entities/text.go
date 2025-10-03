package entities

import "time"

// Text represents a text content block
type Text struct {
	ID            string
	Slug          string
	Content       string
	PageID        string    // Optional
	PageSlug      string    // Optional
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastUpdatedBy string
}
