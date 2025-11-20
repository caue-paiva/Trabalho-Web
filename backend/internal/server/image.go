package server

import (
	"context"
	"fmt"
	"time"

	"backend/internal/entities"
)

// =======================
// IMAGE OPERATIONS
// =======================

func (s *server) GetImageByID(ctx context.Context, id string) (entities.Image, error) {
	return s.db.GetImageByID(ctx, id)
}

func (s *server) GetImagesBySlug(ctx context.Context, slug string) ([]entities.Image, error) {
	normalized := normalizeSlug(slug)
	return s.db.GetImagesBySlug(ctx, normalized)
}

func (s *server) ListAllImages(ctx context.Context) ([]entities.Image, error) {
	return s.db.ListAllImages(ctx)
}

func (s *server) UploadImage(ctx context.Context, meta entities.Image, data []byte) (entities.Image, error) {
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

func (s *server) UpdateImage(ctx context.Context, id string, meta entities.Image, data []byte) (entities.Image, error) {
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

func (s *server) DeleteImage(ctx context.Context, id string) error {
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
