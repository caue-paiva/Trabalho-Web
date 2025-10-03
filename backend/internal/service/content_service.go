package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"backend/internal/entities"
)

// ContentService defines business logic for content operations
type ContentService interface {
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
}

// contentService implements ContentService
type contentService struct {
	db  DBPort
	obj ObjectStorePort
}

// NewContentService creates a new ContentService
func NewContentService(db DBPort, obj ObjectStorePort) ContentService {
	return &contentService{
		db:  db,
		obj: obj,
	}
}

// Text operations implementation

func (s *contentService) GetTextBySlug(ctx context.Context, slug string) (entities.Text, error) {
	normalized := normalizeSlug(slug)
	return s.db.GetTextBySlug(ctx, normalized)
}

func (s *contentService) GetTextByID(ctx context.Context, id string) (entities.Text, error) {
	return s.db.GetTextByID(ctx, id)
}

func (s *contentService) GetTextsByPageID(ctx context.Context, pageID string) ([]entities.Text, error) {
	return s.db.GetTextsByPageID(ctx, pageID)
}

func (s *contentService) GetTextsByPageSlug(ctx context.Context, pageSlug string) ([]entities.Text, error) {
	normalized := normalizeSlug(pageSlug)
	return s.db.ListTextsByPageSlug(ctx, normalized)
}

func (s *contentService) ListAllTexts(ctx context.Context) ([]entities.Text, error) {
	return s.db.ListAllTexts(ctx)
}

func (s *contentService) CreateText(ctx context.Context, text entities.Text) (entities.Text, error) {
	// Business logic: normalize slug
	text.Slug = normalizeSlug(text.Slug)

	// Set audit fields
	now := time.Now()
	text.CreatedAt = now
	text.UpdatedAt = now

	// Delegate to port
	return s.db.CreateText(ctx, text)
}

func (s *contentService) UpdateText(ctx context.Context, id string, text entities.Text) (entities.Text, error) {
	// Set audit fields
	text.UpdatedAt = time.Now()

	// Delegate to port
	return s.db.UpdateText(ctx, id, text)
}

func (s *contentService) DeleteText(ctx context.Context, id string) error {
	return s.db.DeleteText(ctx, id)
}

// Image operations implementation

func (s *contentService) GetImageByID(ctx context.Context, id string) (entities.Image, error) {
	return s.db.GetImageByID(ctx, id)
}

func (s *contentService) GetImagesByGallerySlug(ctx context.Context, slug string) ([]entities.Image, error) {
	normalized := normalizeSlug(slug)
	return s.db.GetImagesByGallerySlug(ctx, normalized)
}

func (s *contentService) UploadImage(ctx context.Context, meta entities.Image, data []byte) (entities.Image, error) {
	// Business logic: generate object key with timestamp
	key := generateObjectKey(meta.Slug)

	// Validate image size (10MB limit)
	if len(data) > 10*1024*1024 {
		return entities.Image{}, fmt.Errorf("image too large: max 10MB")
	}

	// Upload to object store
	url, err := s.obj.PutObject(ctx, key, data)
	if err != nil {
		return entities.Image{}, fmt.Errorf("upload failed: %w", err)
	}

	// Update entity with storage URL and audit fields
	meta.ObjectURL = url
	now := time.Now()
	meta.CreatedAt = now
	meta.UpdatedAt = now

	// Persist metadata
	created, err := s.db.CreateImageMeta(ctx, meta)
	if err != nil {
		// Rollback: delete uploaded object
		_ = s.obj.DeleteObject(ctx, key)
		return entities.Image{}, fmt.Errorf("db persist failed: %w", err)
	}

	return created, nil
}

func (s *contentService) UpdateImage(ctx context.Context, id string, meta entities.Image, data []byte) (entities.Image, error) {
	// If new image data provided, upload it
	if len(data) > 0 {
		// Validate size
		if len(data) > 10*1024*1024 {
			return entities.Image{}, fmt.Errorf("image too large: max 10MB")
		}

		// Generate new key
		key := generateObjectKey(meta.Slug)

		// Upload new image
		url, err := s.obj.PutObject(ctx, key, data)
		if err != nil {
			return entities.Image{}, fmt.Errorf("upload failed: %w", err)
		}

		// Get existing image to delete old object
		existing, err := s.db.GetImageByID(ctx, id)
		if err == nil && existing.ObjectURL != "" {
			// Delete old object (best effort, don't fail if it errors)
			_ = s.obj.DeleteObject(ctx, extractKeyFromURL(existing.ObjectURL))
		}

		meta.ObjectURL = url
	}

	// Set audit fields
	meta.UpdatedAt = time.Now()

	// Update metadata
	return s.db.UpdateImageMeta(ctx, id, meta)
}

func (s *contentService) DeleteImage(ctx context.Context, id string) error {
	// Get image to retrieve object key
	img, err := s.db.GetImageByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete metadata first
	if err := s.db.DeleteImageMeta(ctx, id); err != nil {
		return err
	}

	// Delete object from storage (best effort)
	if img.ObjectURL != "" {
		_ = s.obj.DeleteObject(ctx, extractKeyFromURL(img.ObjectURL))
	}

	return nil
}

// Timeline operations implementation

func (s *contentService) GetTimelineEntryByID(ctx context.Context, id string) (entities.TimelineEntry, error) {
	return s.db.GetTimelineEntryByID(ctx, id)
}

func (s *contentService) ListTimelineEntries(ctx context.Context) ([]entities.TimelineEntry, error) {
	return s.db.ListTimelineEntries(ctx)
}

func (s *contentService) CreateTimelineEntry(ctx context.Context, entry entities.TimelineEntry) (entities.TimelineEntry, error) {
	// Set audit fields
	now := time.Now()
	entry.CreatedAt = now
	entry.UpdatedAt = now

	return s.db.CreateTimelineEntry(ctx, entry)
}

func (s *contentService) UpdateTimelineEntry(ctx context.Context, id string, entry entities.TimelineEntry) (entities.TimelineEntry, error) {
	// Set audit fields
	entry.UpdatedAt = time.Now()

	return s.db.UpdateTimelineEntry(ctx, id, entry)
}

func (s *contentService) DeleteTimelineEntry(ctx context.Context, id string) error {
	return s.db.DeleteTimelineEntry(ctx, id)
}

// Helper functions

func normalizeSlug(slug string) string {
	// Lowercase, trim, replace spaces with hyphens
	normalized := strings.TrimSpace(strings.ToLower(slug))
	normalized = strings.ReplaceAll(normalized, " ", "-")
	return normalized
}

func generateObjectKey(slug string) string {
	// Format: images/{slug}-{timestamp}.jpg
	return fmt.Sprintf("images/%s-%d.jpg", normalizeSlug(slug), time.Now().Unix())
}

func extractKeyFromURL(url string) string {
	// Simple extraction: assume URL ends with the key
	// Example: https://storage.googleapis.com/bucket/images/sunset-123.jpg -> images/sunset-123.jpg
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return ""
}
