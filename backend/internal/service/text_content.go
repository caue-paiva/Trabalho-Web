package service

import (
	"context"
	"time"

	"backend/internal/entities"
)

// TextContentService defines business logic for text and timeline operations
type TextContentService interface {
	// Text operations
	GetTextBySlug(ctx context.Context, slug string) (entities.Text, error)
	GetTextByID(ctx context.Context, id string) (entities.Text, error)
	GetTextsByPageID(ctx context.Context, pageID string) ([]entities.Text, error)
	GetTextsByPageSlug(ctx context.Context, pageSlug string) ([]entities.Text, error)
	ListAllTexts(ctx context.Context) ([]entities.Text, error)
	CreateText(ctx context.Context, text entities.Text) (entities.Text, error)
	UpdateText(ctx context.Context, id string, text entities.Text) (entities.Text, error)
	DeleteText(ctx context.Context, id string) error

	// Timeline operations
	GetTimelineEntryByID(ctx context.Context, id string) (entities.TimelineEntry, error)
	ListTimelineEntries(ctx context.Context) ([]entities.TimelineEntry, error)
	CreateTimelineEntry(ctx context.Context, entry entities.TimelineEntry) (entities.TimelineEntry, error)
	UpdateTimelineEntry(ctx context.Context, id string, entry entities.TimelineEntry) (entities.TimelineEntry, error)
	DeleteTimelineEntry(ctx context.Context, id string) error
}

// textContentService implements TextContentService
type textContentService struct {
	db DBPort
}

// NewTextContentService creates a new TextContentService
func NewTextContentService(db DBPort) TextContentService {
	return &textContentService{
		db: db,
	}
}

// Text operations implementation

func (s *textContentService) GetTextBySlug(ctx context.Context, slug string) (entities.Text, error) {
	normalized := normalizeSlug(slug)
	return s.db.GetTextBySlug(ctx, normalized)
}

func (s *textContentService) GetTextByID(ctx context.Context, id string) (entities.Text, error) {
	return s.db.GetTextByID(ctx, id)
}

func (s *textContentService) GetTextsByPageID(ctx context.Context, pageID string) ([]entities.Text, error) {
	return s.db.GetTextsByPageID(ctx, pageID)
}

func (s *textContentService) GetTextsByPageSlug(ctx context.Context, pageSlug string) ([]entities.Text, error) {
	normalized := normalizeSlug(pageSlug)
	return s.db.ListTextsByPageSlug(ctx, normalized)
}

func (s *textContentService) ListAllTexts(ctx context.Context) ([]entities.Text, error) {
	return s.db.ListAllTexts(ctx)
}

func (s *textContentService) CreateText(ctx context.Context, text entities.Text) (entities.Text, error) {
	// Business logic: normalize slug
	text.Slug = normalizeSlug(text.Slug)

	// Set audit fields
	now := time.Now()
	text.CreatedAt = now
	text.UpdatedAt = now

	// Delegate to port
	return s.db.CreateText(ctx, text)
}

func (s *textContentService) UpdateText(ctx context.Context, id string, text entities.Text) (entities.Text, error) {
	// Set audit fields
	text.UpdatedAt = time.Now()

	// Delegate to port
	return s.db.UpdateText(ctx, id, text)
}

func (s *textContentService) DeleteText(ctx context.Context, id string) error {
	return s.db.DeleteText(ctx, id)
}

// Timeline operations implementation

func (s *textContentService) GetTimelineEntryByID(ctx context.Context, id string) (entities.TimelineEntry, error) {
	return s.db.GetTimelineEntryByID(ctx, id)
}

func (s *textContentService) ListTimelineEntries(ctx context.Context) ([]entities.TimelineEntry, error) {
	return s.db.ListTimelineEntries(ctx)
}

func (s *textContentService) CreateTimelineEntry(ctx context.Context, entry entities.TimelineEntry) (entities.TimelineEntry, error) {
	// Set audit fields
	now := time.Now()
	entry.CreatedAt = now
	entry.UpdatedAt = now

	return s.db.CreateTimelineEntry(ctx, entry)
}

func (s *textContentService) UpdateTimelineEntry(ctx context.Context, id string, entry entities.TimelineEntry) (entities.TimelineEntry, error) {
	// Set audit fields
	entry.UpdatedAt = time.Now()

	return s.db.UpdateTimelineEntry(ctx, id, entry)
}

func (s *textContentService) DeleteTimelineEntry(ctx context.Context, id string) error {
	return s.db.DeleteTimelineEntry(ctx, id)
}
