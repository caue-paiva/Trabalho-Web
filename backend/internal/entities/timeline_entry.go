package entities

import "time"

// TimelineEntry represents a timeline event
type TimelineEntry struct {
	ID            string    `json:"id" firestore:"-"` // Document ID is stored separately, not in document data
	Name          string    `json:"name" firestore:"name"`
	Text          string    `json:"text" firestore:"text"`
	Location      string    `json:"location,omitempty" firestore:"location,omitempty"`
	Date          time.Time `json:"date" firestore:"date"`
	CreatedAt     time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" firestore:"updatedAt"`
	LastUpdatedBy string    `json:"lastUpdatedBy,omitempty" firestore:"lastUpdatedBy,omitempty"`
}
