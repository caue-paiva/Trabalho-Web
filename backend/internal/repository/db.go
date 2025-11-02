package clients

import (
	"context"
	"fmt"

	"backend/internal/entities"
	"backend/internal/service"
)

// Compile-time interface check
var _ service.DBPort = (*dbClient)(nil)

type dbClient struct {
	// Future: repository dependencies
}

// NewDBClient creates a new DBPort implementation (stub)
func NewDBClient() service.DBPort {
	return &dbClient{}
}

// Text operations (stubs)

func (c *dbClient) GetTextBySlug(ctx context.Context, slug string) (entities.Text, error) {
	return entities.Text{}, fmt.Errorf("not implemented")
}

func (c *dbClient) GetTextByID(ctx context.Context, id string) (entities.Text, error) {
	return entities.Text{}, fmt.Errorf("not implemented")
}

func (c *dbClient) GetTextsByPageID(ctx context.Context, pageID string) ([]entities.Text, error) {
	return []entities.Text{}, fmt.Errorf("not implemented")
}

func (c *dbClient) ListTextsByPageSlug(ctx context.Context, pageSlug string) ([]entities.Text, error) {
	return []entities.Text{}, fmt.Errorf("not implemented")
}

func (c *dbClient) ListAllTexts(ctx context.Context) ([]entities.Text, error) {
	return []entities.Text{}, fmt.Errorf("not implemented")
}

func (c *dbClient) CreateText(ctx context.Context, text entities.Text) (entities.Text, error) {
	return entities.Text{}, fmt.Errorf("not implemented")
}

func (c *dbClient) UpdateText(ctx context.Context, id string, patch entities.Text) (entities.Text, error) {
	return entities.Text{}, fmt.Errorf("not implemented")
}

func (c *dbClient) DeleteText(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}

// Image operations (stubs)

func (c *dbClient) GetImageByID(ctx context.Context, id string) (entities.Image, error) {
	return entities.Image{}, fmt.Errorf("not implemented")
}

func (c *dbClient) GetImagesByGallerySlug(ctx context.Context, slug string) ([]entities.Image, error) {
	return []entities.Image{}, fmt.Errorf("not implemented")
}

func (c *dbClient) CreateImageMeta(ctx context.Context, img entities.Image) (entities.Image, error) {
	return entities.Image{}, fmt.Errorf("not implemented")
}

func (c *dbClient) UpdateImageMeta(ctx context.Context, id string, patch entities.Image) (entities.Image, error) {
	return entities.Image{}, fmt.Errorf("not implemented")
}

func (c *dbClient) DeleteImageMeta(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}

// Timeline operations (stubs)

func (c *dbClient) GetTimelineEntryByID(ctx context.Context, id string) (entities.TimelineEntry, error) {
	return entities.TimelineEntry{}, fmt.Errorf("not implemented")
}

func (c *dbClient) ListTimelineEntries(ctx context.Context) ([]entities.TimelineEntry, error) {
	return []entities.TimelineEntry{}, fmt.Errorf("not implemented")
}

func (c *dbClient) CreateTimelineEntry(ctx context.Context, entry entities.TimelineEntry) (entities.TimelineEntry, error) {
	return entities.TimelineEntry{}, fmt.Errorf("not implemented")
}

func (c *dbClient) UpdateTimelineEntry(ctx context.Context, id string, patch entities.TimelineEntry) (entities.TimelineEntry, error) {
	return entities.TimelineEntry{}, fmt.Errorf("not implemented")
}

func (c *dbClient) DeleteTimelineEntry(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}
