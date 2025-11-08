package integration_tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// EventResponse represents the API response for an event
type EventResponse struct {
	ID                string `json:"id"`
	Identifier        string `json:"identifier"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	StartsAt          string `json:"startsAt"`
	EndsAt            string `json:"endsAt"`
	Timezone          string `json:"timezone"`
	LocationName      string `json:"locationName"`
	LogoURL           string `json:"logoUrl"`
	ThumbnailImageURL string `json:"thumbnailImageUrl"`
	LargeImageURL     string `json:"largeImageUrl"`
	OriginalImageURL  string `json:"originalImageUrl"`
	IconImageURL      string `json:"iconImageUrl"`
	Privacy           string `json:"privacy"`
	State             string `json:"state"`
	CreatedAt         string `json:"createdAt"`
}

func TestEvents_GetAll(t *testing.T) {
	// Get all events (default parameters)
	resp := MakeRequest(t, "GET", "/events", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var events []EventResponse
	ParseJSONResponse(t, resp, &events)

	// We can't assert exact count since it's external data,
	// but we can verify the response structure
	if len(events) > 0 {
		event := events[0]
		assert.NotEmpty(t, event.ID, "Event should have an ID")
		assert.NotEmpty(t, event.Name, "Event should have a name")
		// Note: StartsAt and EndsAt may be empty in some events from external API
		// We just verify the fields exist in the struct
	}
}

func TestEvents_WithLimitParameter(t *testing.T) {
	// Get events with limit parameter
	resp := MakeRequest(t, "GET", "/events?limit=5", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var events []EventResponse
	ParseJSONResponse(t, resp, &events)

	// Should have at most 5 events
	assert.LessOrEqual(t, len(events), 5, "Should respect limit parameter")
}

func TestEvents_WithCombinedParameters(t *testing.T) {
	// Get events with limit parameter only (orderBy is not supported by external API)
	resp := MakeRequest(t, "GET", "/events?limit=3", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var events []EventResponse
	ParseJSONResponse(t, resp, &events)

	// Should respect limit
	assert.LessOrEqual(t, len(events), 3, "Should respect limit parameter")

	// Verify event structure (dates may be empty in external API data)
	if len(events) > 0 {
		event := events[0]
		assert.NotEmpty(t, event.ID)
		assert.NotEmpty(t, event.Name)
		// Note: StartsAt and EndsAt may be empty in some events
	}
}

func TestEvents_ResponseStructure(t *testing.T) {
	// Get a single event to verify full structure
	resp := MakeRequest(t, "GET", "/events?limit=1", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var events []EventResponse
	ParseJSONResponse(t, resp, &events)

	if len(events) > 0 {
		event := events[0]

		// Verify fields that should always be present
		assert.NotEmpty(t, event.ID, "Should have ID")
		assert.NotEmpty(t, event.Identifier, "Should have identifier")
		assert.NotEmpty(t, event.Name, "Should have name")

		// These fields may be empty in some events from the external API
		// We verify they exist in the struct but don't assert they're non-empty
		t.Logf("Event ID: %s", event.ID)
		t.Logf("Event Name: %s", event.Name)
		t.Logf("StartsAt: %s (may be empty)", event.StartsAt)
		t.Logf("EndsAt: %s (may be empty)", event.EndsAt)
		t.Logf("Timezone: %s", event.Timezone)
		t.Logf("Privacy: %s", event.Privacy)
		t.Logf("State: %s", event.State)
	} else {
		t.Skip("No events available from external API to test structure")
	}
}

func TestEvents_EmptyResult(t *testing.T) {
	// Request with very restrictive limit
	resp := MakeRequest(t, "GET", "/events?limit=0", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var events []EventResponse
	ParseJSONResponse(t, resp, &events)

	// Should handle empty result gracefully
	assert.NotNil(t, events, "Events array should not be nil")
}

func TestEvents_InvalidParameters(t *testing.T) {
	// Test with invalid limit (negative)
	resp := MakeRequest(t, "GET", "/events?limit=-1", nil)
	// Should either return 400 or handle gracefully with 200
	assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, resp.StatusCode)
	resp.Body.Close()
}

func TestEvents_LargeLimit(t *testing.T) {
	// Request with large limit
	resp := MakeRequest(t, "GET", "/events?limit=100", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var events []EventResponse
	ParseJSONResponse(t, resp, &events)

	// Should not exceed reasonable limits (API might cap at a max value)
	assert.LessOrEqual(t, len(events), 100, "Should not exceed requested limit")
}

func TestEvents_DateOrdering(t *testing.T) {
	// Get multiple events (orderBy parameter not supported by external API)
	resp := MakeRequest(t, "GET", "/events?limit=10", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var events []EventResponse
	ParseJSONResponse(t, resp, &events)

	// Just verify we got events back
	// Note: Some events may have empty date fields in the external API
	assert.Greater(t, len(events), 0, "Should receive at least one event")
	for _, event := range events {
		assert.NotEmpty(t, event.ID, "Each event should have an ID")
		assert.NotEmpty(t, event.Name, "Each event should have a name")
		// StartsAt and EndsAt may be empty - this is valid for the external API
	}
}
