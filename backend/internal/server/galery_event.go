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

// CreateGaleryEvent uploads images to object storage, creates image documents, and creates a galery event
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

	// Upload all images to object storage and create image documents
	imageURLs := make([]string, 0, len(imagesBase64))
	createdImageIDs := make([]string, 0, len(imagesBase64)) // Track created image IDs for rollback

	for i, base64Image := range imagesBase64 {
		// Decode base64 image
		imageData, err := base64.StdEncoding.DecodeString(base64Image)
		if err != nil {
			// Rollback: delete all previously uploaded images and image documents
			s.rollbackGaleryEventCreation(ctx, createdImageIDs)
			return entities.GaleryEvent{}, fmt.Errorf("failed to decode image %d: %w", i, err)
		}

		// Generate unique key for image in object storage
		imageSlug := fmt.Sprintf("galery_events/%s/%s_%d", uuid.New().String(), time.Now().Format("20060102"), i)
		// Create an Image document in Firestore for this photo
		imageMeta := entities.Image{
			Slug:     imageSlug,
			Name:     fmt.Sprintf("%s - Foto %d", name, i+1),
			Text:     fmt.Sprintf("Imagem do evento: %s", name),
			Date:     date,
			Location: location,
		}

		createdImage, err := s.UploadImage(ctx, imageMeta, imageData)
		if err != nil {
			// Rollback: delete all previously uploaded images and image documents
			s.rollbackGaleryEventCreation(ctx, createdImageIDs)
			return entities.GaleryEvent{}, fmt.Errorf("failed to create image document %d: %w", i, err)
		}

		imageURLs = append(imageURLs, createdImage.ObjectURL)
		createdImageIDs = append(createdImageIDs, createdImage.ID)
	}

	// Create galery event entity
	galeryEvent := entities.GaleryEvent{
		Name:      name,
		Location:  location,
		Date:      date,
		ImageURLs: imageURLs,
		ImageIDs:  createdImageIDs,
	}

	// Save galery event to database
	savedEvent, err := s.db.CreateGaleryEvent(ctx, galeryEvent)
	if err != nil {
		// Rollback: delete all uploaded images and image documents
		s.rollbackGaleryEventCreation(ctx, createdImageIDs)
		return entities.GaleryEvent{}, fmt.Errorf("failed to save galery event to database: %w", err)
	}

	return savedEvent, nil
}

// rollbackGaleryEventCreation deletes uploaded images and image documents in case of failure
func (s *server) rollbackGaleryEventCreation(ctx context.Context, imageIDs []string) {
	// Delete created image documents from Firestore
	for _, imageID := range imageIDs {
		// Best effort deletion - log errors but don't fail
		if err := s.db.DeleteImageMeta(ctx, imageID); err != nil {
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

func (s *server) ModifyGaleryEvent(ctx context.Context, id string, newEvent entities.GaleryEvent) (entities.GaleryEvent, error) {
	return s.db.ModifyGaleryEvent(ctx, id, newEvent)
}
