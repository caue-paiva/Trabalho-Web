package entities

import "time"

// Text represents a text content block
type Text struct {
	ID            string    `json:"id" firestore:"-"` // Document ID is stored separately, not in document data
	Slug          string    `json:"slug" firestore:"slug"`
	Content       string    `json:"content" firestore:"content"`
	PageID        string    `json:"pageId,omitempty" firestore:"pageId,omitempty"` // Optional
	PageSlug      string    `json:"pageSlug,omitempty" firestore:"pageSlug,omitempty"` // Optional
	CreatedAt     time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" firestore:"updatedAt"`
	LastUpdatedBy string    `json:"lastUpdatedBy,omitempty" firestore:"lastUpdatedBy,omitempty"`
}
