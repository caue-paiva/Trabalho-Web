package server

import (
	"context"
	"time"

	"backend/internal/entities"
)

// Server defines the unified service interface for all business operations
type Server interface {
	// Text operations
	GetTextBySlug(ctx context.Context, slug string) (entities.Text, error)
	GetTextByID(ctx context.Context, id string) (entities.Text, error)
	GetTextsByPageID(ctx context.Context, pageID string) ([]entities.Text, error)
	GetTextsByPageSlug(ctx context.Context, pageSlug string) ([]entities.Text, error)
	ListAllTexts(ctx context.Context) ([]entities.Text, error)
	CreateText(ctx context.Context, text entities.Text) (entities.Text, error)
	UpdateText(ctx context.Context, id string, text entities.Text) (entities.Text, error)
	DeleteText(ctx context.Context, id string) error

	// Image operations
	GetImageByID(ctx context.Context, id string) (entities.Image, error)
	GetImagesByGallerySlug(ctx context.Context, slug string) ([]entities.Image, error)
	UploadImage(ctx context.Context, meta entities.Image, data []byte) (entities.Image, error)
	UpdateImage(ctx context.Context, id string, meta entities.Image, data []byte) (entities.Image, error)
	DeleteImage(ctx context.Context, id string) error

	// Timeline operations
	GetTimelineEntryByID(ctx context.Context, id string) (entities.TimelineEntry, error)
	ListTimelineEntries(ctx context.Context) ([]entities.TimelineEntry, error)
	CreateTimelineEntry(ctx context.Context, entry entities.TimelineEntry) (entities.TimelineEntry, error)
	UpdateTimelineEntry(ctx context.Context, id string, entry entities.TimelineEntry) (entities.TimelineEntry, error)
	DeleteTimelineEntry(ctx context.Context, id string) error

	// Events operations
	GetEvents(ctx context.Context, limit int, orderBy string, desc bool) ([]entities.Event, error)

	// GaleryEvent operations
	CreateGaleryEvent(ctx context.Context, name, location string, date time.Time, imagesBase64 []string) (entities.GaleryEvent, error)
	GetGaleryEventByID(ctx context.Context, id string) (entities.GaleryEvent, error)
	ListGaleryEvents(ctx context.Context) ([]entities.GaleryEvent, error)
}

// server implements the Server interface
type server struct {
	db     DBPort
	obj    ObjectStorePort
	events GrupyEventsPort
}

// NewServer creates a new unified Server with all dependencies
func NewServer(db DBPort, obj ObjectStorePort, events GrupyEventsPort) Server {
	return &server{
		db:     db,
		obj:    obj,
		events: events,
	}
}
