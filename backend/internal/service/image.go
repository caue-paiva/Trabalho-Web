package service

import (
	"context"
	"fmt"
	"time"

	"backend/internal/entities"
)

// ImageService defines business logic for image operations
type ImageService interface {
	GetImageByID(ctx context.Context, id string) (entities.Image, error)
	GetImagesByGallerySlug(ctx context.Context, slug string) ([]entities.Image, error)
	UploadImage(ctx context.Context, meta entities.Image, data []byte) (entities.Image, error)
	UpdateImage(ctx context.Context, id string, meta entities.Image, data []byte) (entities.Image, error)
	DeleteImage(ctx context.Context, id string) error
}

// imageService implements ImageService
type imageService struct {
	db  DBPort
	obj ObjectStorePort
}

// NewImageService creates a new ImageService
func NewImageService(db DBPort, obj ObjectStorePort) ImageService {
	return &imageService{
		db:  db,
		obj: obj,
	}
}

// GetImageByID retrieves an image by ID
func (s *imageService) GetImageByID(ctx context.Context, id string) (entities.Image, error) {
	return s.db.GetImageByID(ctx, id)
}

// GetImagesByGallerySlug retrieves all images for a gallery
func (s *imageService) GetImagesByGallerySlug(ctx context.Context, slug string) ([]entities.Image, error) {
	normalized := normalizeSlug(slug)
	return s.db.GetImagesByGallerySlug(ctx, normalized)
}

// UploadImage uploads a new image to object storage and persists metadata
func (s *imageService) UploadImage(ctx context.Context, meta entities.Image, data []byte) (entities.Image, error) {
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

// UpdateImage updates image metadata and optionally uploads new image data
func (s *imageService) UpdateImage(ctx context.Context, id string, meta entities.Image, data []byte) (entities.Image, error) {
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

// DeleteImage deletes an image and its object storage data
func (s *imageService) DeleteImage(ctx context.Context, id string) error {
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
