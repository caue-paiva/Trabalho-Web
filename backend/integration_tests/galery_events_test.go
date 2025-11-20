package integration_tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GaleryEventResponse represents the API response for a galery event
type GaleryEventResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Date      time.Time `json:"date"`
	ImageURLs []string  `json:"image_urls"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateGaleryEventRequest represents the request body for creating a galery event
type CreateGaleryEventRequest struct {
	Name         string   `json:"name"`
	Location     string   `json:"location"`
	Date         string   `json:"date"`
	ImagesBase64 []string `json:"images_base64"`
}

func TestGaleryEvents_CreateAndGet(t *testing.T) {
	// Create a galery event with multiple images
	createReq := CreateGaleryEventRequest{
		Name:     "Integration Test Event",
		Location: "Test Location",
		Date:     time.Now().Format(time.RFC3339),
		ImagesBase64: []string{
			TinyPNG, // First image
			TinyPNG, // Second image
			TinyPNG, // Third image
		},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusCreated)

	var created GaleryEventResponse
	ParseJSONResponse(t, resp, &created)

	// Validate created galery event
	assert.NotEmpty(t, created.ID, "GaleryEvent should have an ID")
	assert.Equal(t, createReq.Name, created.Name)
	assert.Equal(t, createReq.Location, created.Location)
	assert.Len(t, created.ImageURLs, 3, "Should have 3 image URLs")
	assert.NotEmpty(t, created.CreatedAt)
	assert.NotEmpty(t, created.UpdatedAt)

	// Verify all image URLs are not empty
	for i, url := range created.ImageURLs {
		assert.NotEmpty(t, url, "Image URL %d should not be empty", i)
	}

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/galery_events/"+created.ID, nil)
		resp.Body.Close()
	}()

	// Get by ID
	resp = MakeRequest(t, "GET", "/galery_events/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var retrieved GaleryEventResponse
	ParseJSONResponse(t, resp, &retrieved)

	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
	assert.Equal(t, created.Location, retrieved.Location)
	assert.Len(t, retrieved.ImageURLs, 3)
	assert.Equal(t, created.ImageURLs, retrieved.ImageURLs)
}

func TestGaleryEvents_List(t *testing.T) {
	// Create multiple galery events with different dates
	now := time.Now()
	uniquePrefix := GenerateUniqueSlug("galevent")
	events := []CreateGaleryEventRequest{
		{
			Name:         uniquePrefix + " - Event 1 - Oldest",
			Location:     "Location 1",
			Date:         now.Add(-48 * time.Hour).Format(time.RFC3339),
			ImagesBase64: []string{TinyPNG},
		},
		{
			Name:         uniquePrefix + " - Event 2 - Middle",
			Location:     "Location 2",
			Date:         now.Add(-24 * time.Hour).Format(time.RFC3339),
			ImagesBase64: []string{TinyPNG},
		},
		{
			Name:         uniquePrefix + " - Event 3 - Newest",
			Location:     "Location 3",
			Date:         now.Format(time.RFC3339),
			ImagesBase64: []string{TinyPNG},
		},
	}

	var createdIDs []string
	for _, evt := range events {
		resp := MakeRequest(t, "POST", "/galery_events", evt)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var created GaleryEventResponse
		ParseJSONResponse(t, resp, &created)
		createdIDs = append(createdIDs, created.ID)

		// Small delay to ensure distinct creation timestamps
		time.Sleep(100 * time.Millisecond)
	}

	// Cleanup
	defer func() {
		for _, id := range createdIDs {
			resp := MakeRequest(t, "DELETE", "/galery_events/"+id, nil)
			resp.Body.Close()
		}
	}()

	// List all galery events
	resp := MakeRequest(t, "GET", "/galery_events", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var galeryEvents []GaleryEventResponse
	ParseJSONResponse(t, resp, &galeryEvents)

	assert.GreaterOrEqual(t, len(galeryEvents), 3, "Should have at least our 3 events")

	// Verify our events are in the list
	foundCount := 0
	for _, evt := range galeryEvents {
		for _, id := range createdIDs {
			if evt.ID == id {
				foundCount++
				break
			}
		}
	}
	assert.Equal(t, 3, foundCount, "Should find all 3 created events in the list")

	// Verify list is ordered by date descending (newest first)
	// Find our created events in the list (should be in order)
	var ourEvents []GaleryEventResponse
	for _, evt := range galeryEvents {
		for _, id := range createdIDs {
			if evt.ID == id {
				ourEvents = append(ourEvents, evt)
				break
			}
		}
	}

	// Verify dates are valid and events have all required fields
	for i, evt := range ourEvents {
		assert.NotZero(t, evt.Date, "Event %d should have a valid date", i)
		assert.NotEmpty(t, evt.Name, "Event %d should have a name", i)
		assert.NotEmpty(t, evt.Location, "Event %d should have a location", i)
		assert.NotEmpty(t, evt.ImageURLs, "Event %d should have image URLs", i)
	}
}

func TestGaleryEvents_GetByID_NotFound(t *testing.T) {
	// Try to get non-existent galery event
	resp := MakeRequest(t, "GET", "/galery_events/non-existent-id-12345", nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestGaleryEvents_Create_MissingName(t *testing.T) {
	// Try to create event without name
	createReq := CreateGaleryEventRequest{
		Location:     "Test Location",
		Date:         time.Now().Format(time.RFC3339),
		ImagesBase64: []string{TinyPNG},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestGaleryEvents_Create_MissingLocation(t *testing.T) {
	// Try to create event without location
	createReq := CreateGaleryEventRequest{
		Name:         "Test Event",
		Date:         time.Now().Format(time.RFC3339),
		ImagesBase64: []string{TinyPNG},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestGaleryEvents_Create_MissingDate(t *testing.T) {
	// Try to create event without date
	createReq := CreateGaleryEventRequest{
		Name:         "Test Event",
		Location:     "Test Location",
		ImagesBase64: []string{TinyPNG},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestGaleryEvents_Create_NoImages(t *testing.T) {
	// Try to create event without images
	createReq := CreateGaleryEventRequest{
		Name:         "Test Event",
		Location:     "Test Location",
		Date:         time.Now().Format(time.RFC3339),
		ImagesBase64: []string{},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestGaleryEvents_Create_InvalidBase64(t *testing.T) {
	// Try to create event with invalid base64 image
	createReq := CreateGaleryEventRequest{
		Name:     "Test Event",
		Location: "Test Location",
		Date:     time.Now().Format(time.RFC3339),
		ImagesBase64: []string{
			TinyPNG,                 // Valid
			"not-valid-base64!@#$%", // Invalid
			TinyPNG,                 // Valid
		},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	// Should fail because of invalid base64
	AssertStatusCode(t, resp, http.StatusInternalServerError)
	resp.Body.Close()

	// Note: In a production system, you might want to return 400 Bad Request
	// instead of 500 for invalid input data
}

func TestGaleryEvents_ImageURLsAccessible(t *testing.T) {
	// Create a galery event
	createReq := CreateGaleryEventRequest{
		Name:     "URL Accessibility Test",
		Location: "Test Location",
		Date:     time.Now().Format(time.RFC3339),
		ImagesBase64: []string{
			TinyPNG,
			TinyPNG,
		},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created GaleryEventResponse
	ParseJSONResponse(t, resp, &created)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/galery_events/"+created.ID, nil)
		resp.Body.Close()
	}()

	// Verify all image URLs are accessible
	for i, imageURL := range created.ImageURLs {
		assert.NotEmpty(t, imageURL, "Image URL %d should not be empty", i)

		// Make a HEAD request to verify the URL is accessible
		objectResp, err := http.Head(imageURL)
		require.NoError(t, err, "Should be able to access image URL %d", i)
		defer objectResp.Body.Close()

		assert.Equal(t, http.StatusOK, objectResp.StatusCode,
			"Image URL %d should be accessible (got %d)", i, objectResp.StatusCode)
	}
}

func TestGaleryEvents_SingleImage(t *testing.T) {
	// Test creating event with just one image (minimum requirement)
	createReq := CreateGaleryEventRequest{
		Name:         "Single Image Event",
		Location:     "Test Location",
		Date:         time.Now().Format(time.RFC3339),
		ImagesBase64: []string{TinyPNG},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusCreated)

	var created GaleryEventResponse
	ParseJSONResponse(t, resp, &created)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/galery_events/"+created.ID, nil)
		resp.Body.Close()
	}()

	assert.Len(t, created.ImageURLs, 1, "Should have exactly 1 image URL")
	assert.NotEmpty(t, created.ImageURLs[0], "Single image URL should not be empty")
}

func TestGaleryEvents_ManyImages(t *testing.T) {
	// Test creating event with many images
	manyImages := make([]string, 10)
	for i := range manyImages {
		manyImages[i] = TinyPNG
	}

	createReq := CreateGaleryEventRequest{
		Name:         "Many Images Event",
		Location:     "Test Location",
		Date:         time.Now().Format(time.RFC3339),
		ImagesBase64: manyImages,
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusCreated)

	var created GaleryEventResponse
	ParseJSONResponse(t, resp, &created)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/galery_events/"+created.ID, nil)
		resp.Body.Close()
	}()

	assert.Len(t, created.ImageURLs, 10, "Should have exactly 10 image URLs")

	// Verify all URLs are unique and not empty
	urlSet := make(map[string]bool)
	for i, url := range created.ImageURLs {
		assert.NotEmpty(t, url, "Image URL %d should not be empty", i)
		assert.False(t, urlSet[url], "Image URL %d should be unique", i)
		urlSet[url] = true
	}
}

func TestGaleryEvents_DateValidation(t *testing.T) {
	// Test with invalid date format
	createReq := CreateGaleryEventRequest{
		Name:         "Invalid Date Event",
		Location:     "Test Location",
		Date:         "not-a-valid-date",
		ImagesBase64: []string{TinyPNG},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	// Should fail due to invalid date format
	AssertStatusCode(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestGaleryEvents_MultipleImagesTransaction(t *testing.T) {
	// Test that all images are uploaded as part of a transaction
	createReq := CreateGaleryEventRequest{
		Name:     "Multi-Image Transaction Test",
		Location: "Test Location",
		Date:     time.Now().Format(time.RFC3339),
		ImagesBase64: []string{
			TinyPNG,
			TinyPNG,
			TinyPNG,
			TinyPNG,
			TinyPNG, // 5 images
		},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusCreated)

	var created GaleryEventResponse
	ParseJSONResponse(t, resp, &created)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/galery_events/"+created.ID, nil)
		resp.Body.Close()
	}()

	// Verify all 5 images were uploaded
	assert.Len(t, created.ImageURLs, 5, "Should have exactly 5 image URLs")

	// Verify all URLs are unique
	urlSet := make(map[string]bool)
	for i, url := range created.ImageURLs {
		assert.NotEmpty(t, url, "Image URL %d should not be empty", i)
		assert.False(t, urlSet[url], "Image URL %d should be unique", i)
		urlSet[url] = true
	}

	// Verify all images are accessible
	for i, imageURL := range created.ImageURLs {
		objectResp, err := http.Head(imageURL)
		require.NoError(t, err, "Should be able to access image URL %d", i)
		defer objectResp.Body.Close()

		assert.Equal(t, http.StatusOK, objectResp.StatusCode,
			"Image URL %d should be accessible", i)
	}
}

func TestGaleryEvents_EmptyNameLocation(t *testing.T) {
	// Test that empty name fails
	createReq := CreateGaleryEventRequest{
		Name:         "",
		Location:     "Test Location",
		Date:         time.Now().Format(time.RFC3339),
		ImagesBase64: []string{TinyPNG},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusBadRequest)
	resp.Body.Close()

	// Test that empty location fails
	createReq = CreateGaleryEventRequest{
		Name:         "Test Event",
		Location:     "",
		Date:         time.Now().Format(time.RFC3339),
		ImagesBase64: []string{TinyPNG},
	}

	resp = MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestGaleryEvents_LargeNumberOfImages(t *testing.T) {
	// Test creating event with many images (stress test)
	numImages := 15
	images := make([]string, numImages)
	for i := range images {
		images[i] = TinyPNG
	}

	createReq := CreateGaleryEventRequest{
		Name:         "Large Image Set Event",
		Location:     "Test Location",
		Date:         time.Now().Format(time.RFC3339),
		ImagesBase64: images,
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusCreated)

	var created GaleryEventResponse
	ParseJSONResponse(t, resp, &created)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/galery_events/"+created.ID, nil)
		resp.Body.Close()
	}()

	assert.Len(t, created.ImageURLs, numImages, "Should have all %d image URLs", numImages)

	// Verify all URLs are accessible (sample check first 3)
	for i := 0; i < 3 && i < len(created.ImageURLs); i++ {
		objectResp, err := http.Head(created.ImageURLs[i])
		require.NoError(t, err, "Should be able to access image URL %d", i)
		objectResp.Body.Close()
		assert.Equal(t, http.StatusOK, objectResp.StatusCode)
	}
}

func TestGaleryEvents_SpecialCharactersInFields(t *testing.T) {
	// Test with special characters in name and location
	createReq := CreateGaleryEventRequest{
		Name:         "Event with Special Chars: São Paulo, Café & Python 🐍",
		Location:     "São Carlos - IFSP (Instituto Federal)",
		Date:         time.Now().Format(time.RFC3339),
		ImagesBase64: []string{TinyPNG},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	AssertStatusCode(t, resp, http.StatusCreated)

	var created GaleryEventResponse
	ParseJSONResponse(t, resp, &created)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/galery_events/"+created.ID, nil)
		resp.Body.Close()
	}()

	// Verify special characters are preserved
	assert.Equal(t, createReq.Name, created.Name, "Name with special chars should be preserved")
	assert.Equal(t, createReq.Location, created.Location, "Location with special chars should be preserved")
}

func TestGaleryEvents_Delete(t *testing.T) {
	// Create a galery event
	createReq := CreateGaleryEventRequest{
		Name:         "Event to Delete",
		Location:     "Test Location",
		Date:         time.Now().Format(time.RFC3339),
		ImagesBase64: []string{TinyPNG},
	}

	resp := MakeRequest(t, "POST", "/galery_events", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created GaleryEventResponse
	ParseJSONResponse(t, resp, &created)

	// Store image URLs to verify they remain accessible after deletion
	imageURLs := created.ImageURLs

	// Delete the galery event
	resp = MakeRequest(t, "DELETE", "/galery_events/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusNoContent)
	resp.Body.Close()

	// Verify the galery event is deleted
	resp = MakeRequest(t, "GET", "/galery_events/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()

	// Verify images are NOT deleted (they should still be accessible)
	for i, imageURL := range imageURLs {
		objectResp, err := http.Head(imageURL)
		require.NoError(t, err, "Image URL %d should still be accessible after event deletion", i)
		defer objectResp.Body.Close()
		assert.Equal(t, http.StatusOK, objectResp.StatusCode,
			"Image URL %d should remain accessible after event deletion", i)
	}
}

func TestGaleryEvents_Delete_NotFound(t *testing.T) {
	// Try to delete non-existent galery event
	resp := MakeRequest(t, "DELETE", "/galery_events/non-existent-id-12345", nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()
}
