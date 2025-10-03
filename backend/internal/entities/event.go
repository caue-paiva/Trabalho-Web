package entities

import "time"

// Event represents an external community event from Grupy Sanca API (proxy only, no persistence)
type Event struct {
	ID                string
	Identifier        string    // Event's unique string identifier (e.g., "b8324ae2")
	Name              string
	Description       string
	StartsAt          time.Time // Event start date/time
	EndsAt            time.Time // Event end date/time
	Timezone          string
	LocationName      string
	LogoURL           string
	ThumbnailImageURL string
	LargeImageURL     string
	OriginalImageURL  string
	IconImageURL      string
	Privacy           string // e.g., "public"
	State             string // e.g., "draft", "published"
	CreatedAt         time.Time
}
