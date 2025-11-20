package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"backend/internal/entities"
	"github.com/google/uuid"
)

// =======================
// GALERY EVENT OPERATIONS
// =======================

// CreateGaleryEvent uploads images to object storage and creates a galery event
// This method is transactional: if any image upload fails, the entire operation fails
func (s *server) CreateGaleryEvent(ctx context.Context, name, location string, date time.Time, imagesBase64 []string) (entities.GaleryEvent, error) {
	// Validate inputs
	if name == "" {
		return entities.GaleryEvent{}, fmt.Errorf("name is required")
	}
	if location == "" {
		return entities.GaleryEvent{}, fmt.Errorf("location is required")
	}
	if date.IsZero() {
		return entities.GaleryEvent{}, fmt.Errorf("date is required")
	}
	if len(imagesBase64) == 0 {
		return entities.GaleryEvent{}, fmt.Errorf("at least one image is required")
	}

	// Upload all images to object storage
	imageURLs := make([]string, 0, len(imagesBase64))
	uploadedKeys := make([]string, 0, len(imagesBase64)) // Track uploaded keys for rollback

	for i, base64Image := range imagesBase64 {
		// Decode base64 image
		imageData, err := base64.StdEncoding.DecodeString(base64Image)
		if err != nil {
			// Rollback: delete all previously uploaded images
			s.rollbackImageUploads(ctx, uploadedKeys)
			return entities.GaleryEvent{}, fmt.Errorf("failed to decode image %d: %w", i, err)
		}

		// Generate unique key for image
		imageKey := fmt.Sprintf("galery_events/%s/%s_%d", uuid.New().String(), time.Now().Format("20060102"), i)

		// Upload to object storage
		publicURL, err := s.obj.PutObject(ctx, imageKey, imageData)
		if err != nil {
			// Rollback: delete all previously uploaded images
			s.rollbackImageUploads(ctx, uploadedKeys)
			return entities.GaleryEvent{}, fmt.Errorf("failed to upload image %d: %w", i, err)
		}

		imageURLs = append(imageURLs, publicURL)
		uploadedKeys = append(uploadedKeys, imageKey)
	}

	// Create galery event entity
	galeryEvent := entities.GaleryEvent{
		Name:      name,
		Location:  location,
		Date:      date,
		ImageURLs: imageURLs,
	}

	// Save to database
	savedEvent, err := s.db.CreateGaleryEvent(ctx, galeryEvent)
	if err != nil {
		// Rollback: delete all uploaded images
		s.rollbackImageUploads(ctx, uploadedKeys)
		return entities.GaleryEvent{}, fmt.Errorf("failed to save galery event to database: %w", err)
	}

	return savedEvent, nil
}

// rollbackImageUploads deletes uploaded images in case of failure
func (s *server) rollbackImageUploads(ctx context.Context, keys []string) {
	for _, key := range keys {
		// Best effort deletion - log errors but don't fail
		if err := s.obj.DeleteObject(ctx, key); err != nil {
			// In production, you might want to log this error
			// For now, we silently continue
			_ = err
		}
	}
}

// GetGaleryEventByID retrieves a galery event by ID
func (s *server) GetGaleryEventByID(ctx context.Context, id string) (entities.GaleryEvent, error) {
	return s.db.GetGaleryEventByID(ctx, id)
}

// ListGaleryEvents retrieves all galery events, ordered by date descending
func (s *server) ListGaleryEvents(ctx context.Context) ([]entities.GaleryEvent, error) {
	return s.db.ListGaleryEvents(ctx)
}

// DeleteGaleryEvent deletes a galery event by ID
// Note: This does NOT delete the associated images from object storage
func (s *server) DeleteGaleryEvent(ctx context.Context, id string) error {
	return s.db.DeleteGaleryEvent(ctx, id)
}
