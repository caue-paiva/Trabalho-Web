package server

import (
	"context"
	"time"

	"backend/internal/entities"
)

// =======================
// TEXT OPERATIONS
// =======================

func (s *server) GetTextBySlug(ctx context.Context, slug string) (entities.Text, error) {
	normalized := normalizeSlug(slug)
	return s.db.GetTextBySlug(ctx, normalized)
}

func (s *server) GetTextByID(ctx context.Context, id string) (entities.Text, error) {
	return s.db.GetTextByID(ctx, id)
}

func (s *server) GetTextsByPageID(ctx context.Context, pageID string) ([]entities.Text, error) {
	return s.db.GetTextsByPageID(ctx, pageID)
}

func (s *server) GetTextsByPageSlug(ctx context.Context, pageSlug string) ([]entities.Text, error) {
	normalized := normalizeSlug(pageSlug)
	return s.db.ListTextsByPageSlug(ctx, normalized)
}

func (s *server) ListAllTexts(ctx context.Context) ([]entities.Text, error) {
	return s.db.ListAllTexts(ctx)
}

func (s *server) CreateText(ctx context.Context, text entities.Text) (entities.Text, error) {
	// Business logic: normalize slug
	text.Slug = normalizeSlug(text.Slug)

	// Set audit fields
	now := time.Now()
	text.CreatedAt = now
	text.UpdatedAt = now

	// Delegate to port
	return s.db.CreateText(ctx, text)
}

func (s *server) UpdateText(ctx context.Context, id string, text entities.Text) (entities.Text, error) {
	// Set audit fields
	text.UpdatedAt = time.Now()

	// Delegate to port
	return s.db.UpdateText(ctx, id, text)
}

func (s *server) DeleteText(ctx context.Context, id string) error {
	return s.db.DeleteText(ctx, id)
}
