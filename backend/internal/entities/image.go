package entities

import "time"

// Image represents image metadata
type Image struct {
	ID            string    `json:"id" firestore:"-"` // Document ID is stored separately, not in document data
	Slug          string    `json:"slug,omitempty" firestore:"slug,omitempty"` // Optional
	ObjectURL     string    `json:"objectUrl" firestore:"objectUrl"` // Storage URL
	Name          string    `json:"name" firestore:"name"`
	Text          string    `json:"text" firestore:"text"` // Description
	Date          time.Time `json:"date,omitempty" firestore:"date,omitempty"`
	Location      string    `json:"location,omitempty" firestore:"location,omitempty"`
	CreatedAt     time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" firestore:"updatedAt"`
	LastUpdatedBy string    `json:"lastUpdatedBy,omitempty" firestore:"lastUpdatedBy,omitempty"`
}
