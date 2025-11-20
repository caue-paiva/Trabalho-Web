package entities

import "time"

// GaleryEvent represents a gallery event with associated images
type GaleryEvent struct {
	ID        string    `firestore:"id"`
	Name      string    `firestore:"name"`
	Location  string    `firestore:"location"`
	Date      time.Time `firestore:"date"`
	ImageURLs []string  `firestore:"image_urls"` // URLs from object storage
	CreatedAt time.Time `firestore:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at"`
}
