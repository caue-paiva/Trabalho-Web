package server

import (
	"context"

	"backend/internal/entities"
)

// DBPort defines the contract for database operations
type DBPort interface {
	// Text operations
	GetTextBySlug(ctx context.Context, slug string) (entities.Text, error)
	GetTextByID(ctx context.Context, id string) (entities.Text, error)
	GetTextsByPageID(ctx context.Context, pageID string) ([]entities.Text, error)
	ListTextsByPageSlug(ctx context.Context, pageSlug string) ([]entities.Text, error)
	ListAllTexts(ctx context.Context) ([]entities.Text, error)
	CreateText(ctx context.Context, text entities.Text) (entities.Text, error)
	UpdateText(ctx context.Context, id string, patch entities.Text) (entities.Text, error)
	DeleteText(ctx context.Context, id string) error

	// Image operations
	GetImageByID(ctx context.Context, id string) (entities.Image, error)
	GetImagesByGallerySlug(ctx context.Context, slug string) ([]entities.Image, error)
	ListAllImages(ctx context.Context) ([]entities.Image, error)
	CreateImageMeta(ctx context.Context, img entities.Image) (entities.Image, error)
	UpdateImageMeta(ctx context.Context, id string, patch entities.Image) (entities.Image, error)
	DeleteImageMeta(ctx context.Context, id string) error

	// Timeline operations
	GetTimelineEntryByID(ctx context.Context, id string) (entities.TimelineEntry, error)
	ListTimelineEntries(ctx context.Context) ([]entities.TimelineEntry, error)
	CreateTimelineEntry(ctx context.Context, entry entities.TimelineEntry) (entities.TimelineEntry, error)
	UpdateTimelineEntry(ctx context.Context, id string, patch entities.TimelineEntry) (entities.TimelineEntry, error)
	DeleteTimelineEntry(ctx context.Context, id string) error

	// GaleryEvent operations
	CreateGaleryEvent(ctx context.Context, event entities.GaleryEvent) (entities.GaleryEvent, error)
	GetGaleryEventByID(ctx context.Context, id string) (entities.GaleryEvent, error)
	ListGaleryEvents(ctx context.Context) ([]entities.GaleryEvent, error)
}

// ObjectStorePort defines the contract for object storage operations
type ObjectStorePort interface {
	PutObject(ctx context.Context, key string, data []byte) (publicURL string, err error)
	DeleteObject(ctx context.Context, key string) error
	SignedURL(ctx context.Context, key string) (string, error)
}

// GrupyEventsPort defines the contract for external events API
type GrupyEventsPort interface {
	GetEvents(ctx context.Context, limit int, orderBy string, desc bool) ([]entities.Event, error)
}
