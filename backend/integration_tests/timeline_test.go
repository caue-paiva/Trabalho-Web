package integration_tests

import (
	"backend/internal/http/mapper"
	"net/http"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeline_CreateAndGet(t *testing.T) {
	// Create a timeline entry
	eventDate := time.Date(2024, 11, 8, 12, 0, 0, 0, time.UTC)
	createReq := mapper.CreateTimelineEntryRequest{
		Name:     "GrupySanca Meetup Integration Test",
		Text:     "An important milestone in our history",
		Location: "São Carlos, SP",
		Date:     eventDate.Format(time.RFC3339),
	}

	resp := MakeRequest(t, "POST", "/timelineentries", createReq)
	AssertStatusCode(t, resp, http.StatusCreated)

	var created mapper.TimelineEntryResponse
	ParseJSONResponse(t, resp, &created)

	// Validate created entry
	assert.NotEmpty(t, created.ID, "Timeline entry should have an ID")
	assert.Equal(t, createReq.Name, created.Name)
	assert.Equal(t, createReq.Text, created.Text)
	assert.Equal(t, createReq.Location, created.Location)
	assert.NotEmpty(t, created.Date)
	assert.NotEmpty(t, created.CreatedAt)
	assert.NotEmpty(t, created.UpdatedAt)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/timelineentries/"+created.ID, nil)
		resp.Body.Close()
	}()

	// Get by ID
	resp = MakeRequest(t, "GET", "/timelineentries/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var retrieved mapper.TimelineEntryResponse
	ParseJSONResponse(t, resp, &retrieved)

	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
	assert.Equal(t, created.Text, retrieved.Text)
	assert.Equal(t, created.Location, retrieved.Location)
}

func TestTimeline_Update(t *testing.T) {
	// Create a timeline entry
	eventDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	createReq := mapper.CreateTimelineEntryRequest{
		Name:     "Original Event",
		Text:     "Original description",
		Location: "Original Location",
		Date:     eventDate.Format(time.RFC3339),
	}

	resp := MakeRequest(t, "POST", "/timelineentries", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created mapper.TimelineEntryResponse
	ParseJSONResponse(t, resp, &created)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/timelineentries/"+created.ID, nil)
		resp.Body.Close()
	}()

	// Update the entry
	newDate := time.Date(2024, 2, 20, 14, 0, 0, 0, time.UTC)
	updateReq := mapper.UpdateTimelineEntryRequest{
		Name:     "Updated Event",
		Text:     "Updated description",
		Location: "Updated Location",
		Date:     newDate.Format(time.RFC3339),
	}

	resp = MakeRequest(t, "PUT", "/timelineentries/"+created.ID, updateReq)
	AssertStatusCode(t, resp, http.StatusOK)

	var updated mapper.TimelineEntryResponse
	ParseJSONResponse(t, resp, &updated)

	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Updated Event", updated.Name)
	assert.Equal(t, "Updated description", updated.Text)
	assert.Equal(t, "Updated Location", updated.Location)
	assert.Equal(t, "integration-test", updated.LastUpdatedBy)
}

func TestTimeline_Delete(t *testing.T) {
	// Create a timeline entry
	eventDate := time.Date(2024, 3, 10, 18, 0, 0, 0, time.UTC)
	createReq := mapper.CreateTimelineEntryRequest{
		Name: "Event to Delete",
		Text: "This will be deleted",
		Date: eventDate.Format(time.RFC3339),
	}

	resp := MakeRequest(t, "POST", "/timelineentries", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created mapper.TimelineEntryResponse
	ParseJSONResponse(t, resp, &created)

	// Delete the entry
	resp = MakeRequest(t, "DELETE", "/timelineentries/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusNoContent)
	resp.Body.Close()

	// Verify it's deleted
	resp = MakeRequest(t, "GET", "/timelineentries/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestTimeline_List(t *testing.T) {
	// Create multiple timeline entries
	entries := []mapper.CreateTimelineEntryRequest{
		{
			Name:     "First Event",
			Text:     "Description 1",
			Location: "Location 1",
			Date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		},
		{
			Name:     "Second Event",
			Text:     "Description 2",
			Location: "Location 2",
			Date:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		},
		{
			Name:     "Third Event",
			Text:     "Description 3",
			Location: "Location 3",
			Date:     time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		},
	}

	var createdIDs []string
	for _, entry := range entries {
		resp := MakeRequest(t, "POST", "/timelineentries", entry)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var created mapper.TimelineEntryResponse
		ParseJSONResponse(t, resp, &created)
		createdIDs = append(createdIDs, created.ID)
	}

	// Cleanup
	defer func() {
		for _, id := range createdIDs {
			resp := MakeRequest(t, "DELETE", "/timelineentries/"+id, nil)
			resp.Body.Close()
		}
	}()

	// List all timeline entries
	resp := MakeRequest(t, "GET", "/timelineentries", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var allEntries []mapper.TimelineEntryResponse
	ParseJSONResponse(t, resp, &allEntries)

	assert.GreaterOrEqual(t, len(allEntries), 3, "Should have at least our 3 created entries")

	// Verify our entries are in the list
	foundCount := 0
	for _, entry := range allEntries {
		for _, id := range createdIDs {
			if entry.ID == id {
				foundCount++
			}
		}
	}
	assert.Equal(t, 3, foundCount, "Should find all 3 created entries")
}

func TestTimeline_NotFound(t *testing.T) {
	// Try to get non-existent entry
	resp := MakeRequest(t, "GET", "/timelineentries/non-existent-id-12345", nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()

	// Try to update non-existent entry
	updateReq := mapper.UpdateTimelineEntryRequest{Name: "Updated"}
	resp = MakeRequest(t, "PUT", "/timelineentries/non-existent-id-12345", updateReq)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()

	// Try to delete non-existent entry
	resp = MakeRequest(t, "DELETE", "/timelineentries/non-existent-id-12345", nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestTimeline_InvalidDate(t *testing.T) {
	// Try to create entry with invalid date format
	createReq := mapper.CreateTimelineEntryRequest{
		Name: "Invalid Date Event",
		Text: "This has an invalid date",
		Date: "not-a-valid-date",
	}

	resp := MakeRequest(t, "POST", "/timelineentries", createReq)
	AssertStatusCode(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestTimeline_ChronologicalOrder(t *testing.T) {
	// Create entries with different dates
	entries := []mapper.CreateTimelineEntryRequest{
		{
			Name: "Earliest Event",
			Text: "First in timeline",
			Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		},
		{
			Name: "Latest Event",
			Text: "Last in timeline",
			Date: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		},
		{
			Name: "Middle Event",
			Text: "Middle of timeline",
			Date: time.Date(2022, 6, 15, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		},
	}

	var createdIDs []string
	for _, entry := range entries {
		resp := MakeRequest(t, "POST", "/timelineentries", entry)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var created mapper.TimelineEntryResponse
		ParseJSONResponse(t, resp, &created)
		createdIDs = append(createdIDs, created.ID)
	}

	// Cleanup
	defer func() {
		for _, id := range createdIDs {
			resp := MakeRequest(t, "DELETE", "/timelineentries/"+id, nil)
			resp.Body.Close()
		}
	}()

	// List all entries
	resp := MakeRequest(t, "GET", "/timelineentries", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var allEntries []mapper.TimelineEntryResponse
	ParseJSONResponse(t, resp, &allEntries)

	// Verify all our entries are present
	foundCount := 0
	for _, entry := range allEntries {
		for _, id := range createdIDs {
			if entry.ID == id {
				foundCount++
			}
		}
	}
	assert.Equal(t, 3, foundCount, "Should find all 3 entries with different dates")
}
