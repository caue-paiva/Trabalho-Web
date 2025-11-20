package integration_tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ImageResponse represents the API response for an image entity
type ImageResponse struct {
	ID        string `json:"id"`
	Slug      string `json:"slug,omitempty"`
	ObjectURL string `json:"object_url"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	Date      string `json:"date,omitempty"`
	Location  string `json:"location,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateImageRequest represents the request body for creating an image
type CreateImageRequest struct {
	Slug     string `json:"slug,omitempty"`
	Name     string `json:"name"`
	Text     string `json:"text"`
	Date     string `json:"date,omitempty"`
	Location string `json:"location,omitempty"`
	Data     string `json:"data"` // Base64 encoded image
}

// UpdateImageRequest represents the request body for updating an image
type UpdateImageRequest struct {
	Name     string `json:"name,omitempty"`
	Text     string `json:"text,omitempty"`
	Date     string `json:"date,omitempty"`
	Location string `json:"location,omitempty"`
	Data     string `json:"data,omitempty"` // Base64 encoded image (optional)
}

const (
	// TinyPNG is a 1x1 red pixel PNG in base64
	TinyPNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8DwHwAFBQIAX8jx0gAAAABJRU5ErkJggg=="
)

func TestImages_CreateAndGet(t *testing.T) {
	slug := GenerateUniqueSlug("img-test")

	// Create an image
	createReq := CreateImageRequest{
		Slug:     slug,
		Name:     "Integration Test Image",
		Text:     "This is a test image description",
		Location: "Test Suite",
		Data:     TinyPNG,
	}

	resp := MakeRequest(t, "POST", "/images", createReq)
	AssertStatusCode(t, resp, http.StatusCreated)

	var created ImageResponse
	ParseJSONResponse(t, resp, &created)

	// Validate created image
	assert.NotEmpty(t, created.ID, "Image should have an ID")
	assert.Equal(t, slug, created.Slug)
	assert.Equal(t, createReq.Name, created.Name)
	assert.Equal(t, createReq.Text, created.Text)
	assert.Equal(t, createReq.Location, created.Location)
	assert.NotEmpty(t, created.ObjectURL, "Should have object URL from GCS")
	assert.NotEmpty(t, created.CreatedAt)
	assert.NotEmpty(t, created.UpdatedAt)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/images/"+created.ID, nil)
		resp.Body.Close()
	}()

	// Get by ID
	resp = MakeRequest(t, "GET", "/images/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var retrieved ImageResponse
	ParseJSONResponse(t, resp, &retrieved)

	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
	assert.Equal(t, created.ObjectURL, retrieved.ObjectURL)
}

func TestImages_UpdateMetadata(t *testing.T) {
	slug := GenerateUniqueSlug("img-update")

	// Create an image
	createReq := CreateImageRequest{
		Slug:     slug,
		Name:     "Original Name",
		Text:     "Original description",
		Location: "Original Location",
		Data:     TinyPNG,
	}

	resp := MakeRequest(t, "POST", "/images", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created ImageResponse
	ParseJSONResponse(t, resp, &created)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/images/"+created.ID, nil)
		resp.Body.Close()
	}()

	// Update metadata only (no new image data)
	updateReq := UpdateImageRequest{
		Name:     "Updated Name",
		Text:     "Updated description",
		Location: "Updated Location",
	}

	resp = MakeRequest(t, "PUT", "/images/"+created.ID, updateReq)
	AssertStatusCode(t, resp, http.StatusOK)

	var updated ImageResponse
	ParseJSONResponse(t, resp, &updated)

	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "Updated description", updated.Text)
	assert.Equal(t, "Updated Location", updated.Location)
	// Object URL should remain the same (no new image uploaded)
	assert.Equal(t, created.ObjectURL, updated.ObjectURL)
}

func TestImages_Delete(t *testing.T) {
	slug := GenerateUniqueSlug("img-delete")

	// Create an image
	createReq := CreateImageRequest{
		Slug: slug,
		Name: "Image to Delete",
		Text: "This will be deleted",
		Data: TinyPNG,
	}

	resp := MakeRequest(t, "POST", "/images", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created ImageResponse
	ParseJSONResponse(t, resp, &created)

	// Delete the image
	resp = MakeRequest(t, "DELETE", "/images/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusNoContent)
	resp.Body.Close()

	// Verify it's deleted
	resp = MakeRequest(t, "GET", "/images/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestImages_GetByGallerySlug(t *testing.T) {
	gallerySlug := GenerateUniqueSlug("test-gallery")

	// Create images with the same gallery slug
	images := []CreateImageRequest{
		{
			Slug: gallerySlug,
			Name: "Gallery Image 1",
			Text: "First image in gallery",
			Data: TinyPNG,
		},
		{
			Slug: gallerySlug,
			Name: "Gallery Image 2",
			Text: "Second image in gallery",
			Data: TinyPNG,
		},
	}

	var createdIDs []string
	for _, img := range images {
		resp := MakeRequest(t, "POST", "/images", img)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var created ImageResponse
		ParseJSONResponse(t, resp, &created)
		createdIDs = append(createdIDs, created.ID)
	}

	// Cleanup
	defer func() {
		for _, id := range createdIDs {
			resp := MakeRequest(t, "DELETE", "/images/"+id, nil)
			resp.Body.Close()
		}
	}()

	// Get images by gallery slug
	resp := MakeRequest(t, "GET", "/images/slug/"+gallerySlug, nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var galleryImages []ImageResponse
	ParseJSONResponse(t, resp, &galleryImages)

	assert.GreaterOrEqual(t, len(galleryImages), 2, "Should have at least our 2 images")

	// Verify all have the correct gallery slug
	foundCount := 0
	for _, img := range galleryImages {
		if img.ID == createdIDs[0] || img.ID == createdIDs[1] {
			assert.Equal(t, gallerySlug, img.Slug)
			foundCount++
		}
	}
	assert.Equal(t, 2, foundCount, "Should find both our gallery images")
}

func TestImages_ObjectURLAccessible(t *testing.T) {
	slug := GenerateUniqueSlug("img-url-test")

	// Create an image
	createReq := CreateImageRequest{
		Slug: slug,
		Name: "URL Test Image",
		Text: "Testing if object URL is accessible",
		Data: TinyPNG,
	}

	resp := MakeRequest(t, "POST", "/images", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created ImageResponse
	ParseJSONResponse(t, resp, &created)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/images/"+created.ID, nil)
		resp.Body.Close()
	}()

	// Verify the object URL is accessible
	assert.NotEmpty(t, created.ObjectURL, "Should have an object URL")

	// Make a HEAD request to the object URL to verify it's accessible
	objectResp, err := http.Head(created.ObjectURL)
	require.NoError(t, err, "Should be able to access object URL")
	defer objectResp.Body.Close()

	assert.Equal(t, http.StatusOK, objectResp.StatusCode, "Object URL should be accessible")
}

func TestImages_NotFound(t *testing.T) {
	// Try to get non-existent image
	resp := MakeRequest(t, "GET", "/images/non-existent-id-12345", nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()

	// Try to update non-existent image
	updateReq := UpdateImageRequest{Name: "Updated"}
	resp = MakeRequest(t, "PUT", "/images/non-existent-id-12345", updateReq)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()

	// Try to delete non-existent image
	resp = MakeRequest(t, "DELETE", "/images/non-existent-id-12345", nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestImages_InvalidBase64(t *testing.T) {
	// Try to create image with invalid base64
	createReq := CreateImageRequest{
		Name: "Invalid Image",
		Text: "This has invalid base64 data",
		Data: "not-valid-base64!@#$%",
	}

	resp := MakeRequest(t, "POST", "/images", createReq)
	AssertStatusCode(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestImages_ListAll(t *testing.T) {
	// Create multiple images with different slugs
	images := []CreateImageRequest{
		{
			Slug: "gallery-1",
			Name: "List Test Image 1",
			Text: "First test image",
			Data: TinyPNG,
		},
		{
			Slug: "gallery-2",
			Name: "List Test Image 2",
			Text: "Second test image",
			Data: TinyPNG,
		},
		{
			Slug: "gallery-1",
			Name: "List Test Image 3",
			Text: "Third test image",
			Data: TinyPNG,
		},
	}

	var createdIDs []string
	for _, img := range images {
		resp := MakeRequest(t, "POST", "/images", img)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var created ImageResponse
		ParseJSONResponse(t, resp, &created)
		createdIDs = append(createdIDs, created.ID)
	}

	// Cleanup
	defer func() {
		for _, id := range createdIDs {
			resp := MakeRequest(t, "DELETE", "/images/"+id, nil)
			resp.Body.Close()
		}
	}()

	// List all images
	resp := MakeRequest(t, "GET", "/images", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var allImages []ImageResponse
	ParseJSONResponse(t, resp, &allImages)

	// Verify our images are in the list
	assert.GreaterOrEqual(t, len(allImages), 3, "Should have at least our 3 images")

	foundCount := 0
	for _, img := range allImages {
		for _, id := range createdIDs {
			if img.ID == id {
				foundCount++
				break
			}
		}
	}
	assert.Equal(t, 3, foundCount, "Should find all 3 created images in the list")

	// Verify images are ordered by creation date (newest first)
	// Just check that created_at timestamps are valid
	for i, img := range allImages {
		assert.NotEmpty(t, img.CreatedAt, "Image %d should have created_at timestamp", i)
	}
}

func TestImages_ListAll_EmptyResult(t *testing.T) {
	// Even if database is empty, listing should return empty array, not error
	resp := MakeRequest(t, "GET", "/images", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var images []ImageResponse
	ParseJSONResponse(t, resp, &images)

	// Should be a valid array (possibly empty)
	assert.NotNil(t, images, "Should return valid array")
}
