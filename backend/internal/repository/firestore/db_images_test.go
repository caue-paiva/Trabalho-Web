package firestore

import (
	"context"
	"testing"
	"time"

	"backend/internal/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBRepository_CreateImageMeta(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name         string
		image        entities.Image
		expectError  bool
		validateFunc func(t *testing.T, created entities.Image)
	}{
		{
			name: "create image with all fields",
			image: entities.Image{
				Slug:      "test-gallery-full",
				ObjectURL: "https://storage.example.com/images/test.jpg",
				Name:      "Test Image",
				Text:      "A beautiful test image",
				Location:  "São Carlos, SP",
				Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
			validateFunc: func(t *testing.T, created entities.Image) {
				assert.NotEmpty(t, created.ID, "Should have an ID")
				assert.Equal(t, "test-gallery-full", created.Slug)
				assert.Equal(t, "https://storage.example.com/images/test.jpg", created.ObjectURL)
				assert.Equal(t, "Test Image", created.Name)
				assert.Equal(t, "A beautiful test image", created.Text)
				assert.Equal(t, "São Carlos, SP", created.Location)
				assert.False(t, created.Date.IsZero(), "Date should be set")
			},
		},
		{
			name: "create image with minimal required fields",
			image: entities.Image{
				ObjectURL: "https://storage.example.com/images/minimal.jpg",
				Name:      "Minimal Image",
				Text:      "Image with only required fields",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
			validateFunc: func(t *testing.T, created entities.Image) {
				assert.NotEmpty(t, created.ID, "Should have an ID")
				assert.Equal(t, "Minimal Image", created.Name)
				assert.Equal(t, "https://storage.example.com/images/minimal.jpg", created.ObjectURL)
				assert.Empty(t, created.Slug, "Slug should be empty")
				assert.Empty(t, created.Location, "Location should be empty")
				assert.True(t, created.Date.IsZero(), "Date should be zero time")
			},
		},
		{
			name: "create image without objectUrl",
			image: entities.Image{
				Name:      "Image without URL",
				Text:      "This should still work",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
			validateFunc: func(t *testing.T, created entities.Image) {
				assert.NotEmpty(t, created.ID, "Should have an ID")
				assert.Equal(t, "Image without URL", created.Name)
				assert.Empty(t, created.ObjectURL, "ObjectURL should be empty")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the image
			created, err := db.CreateImageMeta(ctx, tt.image)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err, "Failed to create image")

			// Cleanup
			defer func() {
				err := db.DeleteImageMeta(ctx, created.ID)
				assert.NoError(t, err, "Failed to cleanup created image")
			}()

			// Run validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, created)
			}

			// Verify we can retrieve it
			retrieved, err := db.GetImageByID(ctx, created.ID)
			require.NoError(t, err, "Failed to get created image")
			assert.Equal(t, created.ID, retrieved.ID)
			assert.Equal(t, created.Name, retrieved.Name)
		})
	}
}

func TestDBRepository_GetImageByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test image for successful retrieval
	testImage := entities.Image{
		Slug:      "test-get-by-id",
		ObjectURL: "https://storage.example.com/images/get-test.jpg",
		Name:      "Get Test Image",
		Text:      "Image for get by ID test",
		Location:  "São Carlos, SP",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := db.CreateImageMeta(ctx, testImage)
	require.NoError(t, err, "Failed to create test image")

	defer func() {
		db.DeleteImageMeta(ctx, created.ID)
	}()

	tests := []struct {
		name         string
		imageID      string
		expectError  bool
		errorMsg     string
		validateFunc func(t *testing.T, image entities.Image)
	}{
		{
			name:        "get existing image by ID",
			imageID:     created.ID,
			expectError: false,
			validateFunc: func(t *testing.T, image entities.Image) {
				assert.Equal(t, created.ID, image.ID)
				assert.Equal(t, "Get Test Image", image.Name)
				assert.Equal(t, "test-get-by-id", image.Slug)
				assert.Equal(t, "São Carlos, SP", image.Location)
			},
		},
		{
			name:        "get non-existent image by ID",
			imageID:     "non-existent-image-id-12345",
			expectError: true,
			errorMsg:    "not found",
		},
		{
			name:        "get with empty ID",
			imageID:     "",
			expectError: true,
			errorMsg:    "", // Don't check specific error message for invalid input
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image, err := db.GetImageByID(ctx, tt.imageID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				return
			}

			require.NoError(t, err, "Failed to get image")
			if tt.validateFunc != nil {
				tt.validateFunc(t, image)
			}
		})
	}
}

func TestDBRepository_UpdateImageMeta(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name         string
		setupImage   entities.Image
		updatePatch  entities.Image
		expectError  bool
		errorMsg     string
		validateFunc func(t *testing.T, original, updated entities.Image)
	}{
		{
			name: "full update of all fields",
			setupImage: entities.Image{
				Slug:      "test-full-update",
				ObjectURL: "https://storage.example.com/images/original.jpg",
				Name:      "Original Image",
				Text:      "Original description",
				Location:  "São Carlos, SP",
				Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			updatePatch: entities.Image{
				Name:          "Updated Image",
				Text:          "Updated description",
				Location:      "Campinas, SP",
				ObjectURL:     "https://storage.example.com/images/updated.jpg",
				Date:          time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC),
				LastUpdatedBy: "test-user",
			},
			expectError: false,
			validateFunc: func(t *testing.T, original, updated entities.Image) {
				assert.Equal(t, "Updated Image", updated.Name)
				assert.Equal(t, "Updated description", updated.Text)
				assert.Equal(t, "Campinas, SP", updated.Location)
				assert.Equal(t, "https://storage.example.com/images/updated.jpg", updated.ObjectURL)
				assert.Equal(t, "test-user", updated.LastUpdatedBy)
			},
		},
		{
			name: "partial update - only text and lastUpdatedBy",
			setupImage: entities.Image{
				Slug:      "test-partial-update",
				ObjectURL: "https://storage.example.com/images/test.jpg",
				Name:      "Original Name",
				Text:      "Original Text",
				Location:  "Original Location",
				Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			updatePatch: entities.Image{
				Text:          "Only text updated",
				LastUpdatedBy: "admin",
			},
			expectError: false,
			validateFunc: func(t *testing.T, original, updated entities.Image) {
				assert.Equal(t, "Only text updated", updated.Text)
				assert.Equal(t, "admin", updated.LastUpdatedBy)
				// Other fields should remain unchanged
				assert.Equal(t, original.Name, updated.Name)
				assert.Equal(t, original.ObjectURL, updated.ObjectURL)
				assert.Equal(t, original.Location, updated.Location)
			},
		},
		{
			name: "update name only",
			setupImage: entities.Image{
				Slug:      "test-name-update",
				ObjectURL: "https://storage.example.com/images/test.jpg",
				Name:      "Old Name",
				Text:      "Some text",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			updatePatch: entities.Image{
				Name: "New Name",
			},
			expectError: false,
			validateFunc: func(t *testing.T, original, updated entities.Image) {
				assert.Equal(t, "New Name", updated.Name)
				assert.Equal(t, original.Text, updated.Text)
				assert.Equal(t, original.ObjectURL, updated.ObjectURL)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the initial image
			created, err := db.CreateImageMeta(ctx, tt.setupImage)
			require.NoError(t, err, "Failed to create setup image")

			defer func() {
				db.DeleteImageMeta(ctx, created.ID)
			}()

			// Perform update
			updated, err := db.UpdateImageMeta(ctx, created.ID, tt.updatePatch)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				return
			}

			require.NoError(t, err, "Failed to update image")

			// Run validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, created, updated)
			}

			// Verify by fetching again
			retrieved, err := db.GetImageByID(ctx, created.ID)
			require.NoError(t, err, "Failed to get updated image")
			assert.Equal(t, updated.Name, retrieved.Name)
			assert.Equal(t, updated.Text, retrieved.Text)
		})
	}
}

func TestDBRepository_UpdateImageMeta_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	patch := entities.Image{
		Name: "Updated Name",
		Text: "Updated Text",
	}

	_, err := db.UpdateImageMeta(ctx, "non-existent-image-id-12345", patch)
	assert.Error(t, err, "Should return error when updating non-existent image")
	assert.Contains(t, err.Error(), "not found", "Error should mention 'not found'")
}

func TestDBRepository_GetImagesByGallerySlug(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	gallerySlug := "test-gallery-multiple"

	// Create test images
	testImages := []entities.Image{
		{
			Slug:      gallerySlug,
			ObjectURL: "https://storage.example.com/images/img1.jpg",
			Name:      "Image 1",
			Text:      "First image in gallery",
			Location:  "São Carlos, SP",
			Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Slug:      gallerySlug,
			ObjectURL: "https://storage.example.com/images/img2.jpg",
			Name:      "Image 2",
			Text:      "Second image in gallery",
			Location:  "São Paulo, SP",
			Date:      time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Slug:      "other-gallery",
			ObjectURL: "https://storage.example.com/images/img3.jpg",
			Name:      "Image 3",
			Text:      "Image in different gallery",
			Location:  "Campinas, SP",
			Date:      time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	var createdIDs []string
	for _, image := range testImages {
		created, err := db.CreateImageMeta(ctx, image)
		require.NoError(t, err, "Failed to create image")
		createdIDs = append(createdIDs, created.ID)
	}

	defer func() {
		for _, id := range createdIDs {
			db.DeleteImageMeta(ctx, id)
		}
	}()

	tests := []struct {
		name             string
		slug             string
		expectedMinCount int
		validateFunc     func(t *testing.T, images []entities.Image)
	}{
		{
			name:             "get images from gallery with multiple images",
			slug:             gallerySlug,
			expectedMinCount: 2,
			validateFunc: func(t *testing.T, images []entities.Image) {
				// Verify all returned images have the correct gallery slug
				foundCount := 0
				for _, img := range images {
					for i, id := range createdIDs[:2] {
						if img.ID == id {
							assert.Equal(t, gallerySlug, img.Slug)
							assert.Equal(t, testImages[i].Name, img.Name)
							foundCount++
						}
					}
				}
				assert.Equal(t, 2, foundCount, "Should find both gallery images")
			},
		},
		{
			name:             "get images from gallery with single image",
			slug:             "other-gallery",
			expectedMinCount: 1,
			validateFunc: func(t *testing.T, images []entities.Image) {
				foundCount := 0
				for _, img := range images {
					if img.ID == createdIDs[2] {
						assert.Equal(t, "other-gallery", img.Slug)
						assert.Equal(t, "Image 3", img.Name)
						foundCount++
					}
				}
				assert.Equal(t, 1, foundCount, "Should find the single image")
			},
		},
		{
			name:             "get images from non-existent gallery",
			slug:             "non-existent-gallery-12345",
			expectedMinCount: 0,
			validateFunc: func(t *testing.T, images []entities.Image) {
				assert.Empty(t, images, "Should return empty slice for non-existent gallery")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			images, err := db.GetImagesByGallerySlug(ctx, tt.slug)
			require.NoError(t, err, "Should not return error")
			assert.GreaterOrEqual(t, len(images), tt.expectedMinCount, "Should have expected number of images")

			if tt.validateFunc != nil {
				tt.validateFunc(t, images)
			}
		})
	}
}

func TestDBRepository_DeleteImageMeta(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test image
	newImage := entities.Image{
		Slug:      "test-delete",
		ObjectURL: "https://storage.example.com/images/delete-me.jpg",
		Name:      "Image to Delete",
		Text:      "This image will be deleted",
		Location:  "São Carlos, SP",
		Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := db.CreateImageMeta(ctx, newImage)
	require.NoError(t, err, "Failed to create image")

	// Delete the image
	err = db.DeleteImageMeta(ctx, created.ID)
	require.NoError(t, err, "Failed to delete image")

	// Verify it's deleted by trying to get it
	_, err = db.GetImageByID(ctx, created.ID)
	assert.Error(t, err, "Should get error when fetching deleted image")
	assert.Contains(t, err.Error(), "not found", "Error should mention 'not found'")
}
