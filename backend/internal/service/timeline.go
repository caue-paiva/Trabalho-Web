package service

import (
	"context"
	"time"

	"backend/internal/entities"
)

// =======================
// TIMELINE OPERATIONS
// =======================

func (s *server) GetTimelineEntryByID(ctx context.Context, id string) (entities.TimelineEntry, error) {
	return s.db.GetTimelineEntryByID(ctx, id)
}

func (s *server) ListTimelineEntries(ctx context.Context) ([]entities.TimelineEntry, error) {
	return s.db.ListTimelineEntries(ctx)
}

func (s *server) CreateTimelineEntry(ctx context.Context, entry entities.TimelineEntry) (entities.TimelineEntry, error) {
	// Set audit fields
	now := time.Now()
	entry.CreatedAt = now
	entry.UpdatedAt = now

	return s.db.CreateTimelineEntry(ctx, entry)
}

func (s *server) UpdateTimelineEntry(ctx context.Context, id string, entry entities.TimelineEntry) (entities.TimelineEntry, error) {
	// Set audit fields
	entry.UpdatedAt = time.Now()

	return s.db.UpdateTimelineEntry(ctx, id, entry)
}

func (s *server) DeleteTimelineEntry(ctx context.Context, id string) error {
	return s.db.DeleteTimelineEntry(ctx, id)
}
