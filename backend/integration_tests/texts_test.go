package integration_tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TextResponse represents the API response for a text entity
type TextResponse struct {
	ID            string `json:"id"`
	Slug          string `json:"slug"`
	Content       string `json:"content"`
	PageID        string `json:"pageId,omitempty"`
	PageSlug      string `json:"pageSlug,omitempty"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
}

// CreateTextRequest represents the request body for creating a text
type CreateTextRequest struct {
	Slug     string `json:"slug"`
	Content  string `json:"content"`
	PageSlug string `json:"pageSlug,omitempty"`
}

// UpdateTextRequest represents the request body for updating a text
type UpdateTextRequest struct {
	Content       string `json:"content,omitempty"`
	PageSlug      string `json:"pageSlug,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
}

func TestTexts_CreateAndGet(t *testing.T) {
	slug := GenerateUniqueSlug("integration-test")

	// Create a text
	createReq := CreateTextRequest{
		Slug:     slug,
		Content:  "This is an integration test content",
		PageSlug: "test-page",
	}

	resp := MakeRequest(t, "POST", "/texts", createReq)
	AssertStatusCode(t, resp, http.StatusCreated)

	var created TextResponse
	ParseJSONResponse(t, resp, &created)

	// Validate created text
	assert.NotEmpty(t, created.ID, "Text should have an ID")
	assert.Equal(t, slug, created.Slug)
	assert.Equal(t, createReq.Content, created.Content)
	assert.Equal(t, createReq.PageSlug, created.PageSlug)
	assert.NotEmpty(t, created.CreatedAt)
	assert.NotEmpty(t, created.UpdatedAt)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/texts/"+created.ID, nil)
		resp.Body.Close()
	}()

	// Get by slug
	resp = MakeRequest(t, "GET", "/texts/"+slug, nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var retrieved TextResponse
	ParseJSONResponse(t, resp, &retrieved)

	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Slug, retrieved.Slug)
	assert.Equal(t, created.Content, retrieved.Content)

	// Get by ID
	resp = MakeRequest(t, "GET", "/texts/id/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var retrievedByID TextResponse
	ParseJSONResponse(t, resp, &retrievedByID)

	assert.Equal(t, created.ID, retrievedByID.ID)
	assert.Equal(t, created.Slug, retrievedByID.Slug)
}

func TestTexts_Update(t *testing.T) {
	slug := GenerateUniqueSlug("update-test")

	// Create a text
	createReq := CreateTextRequest{
		Slug:     slug,
		Content:  "Original content",
		PageSlug: "original-page",
	}

	resp := MakeRequest(t, "POST", "/texts", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created TextResponse
	ParseJSONResponse(t, resp, &created)

	// Cleanup
	defer func() {
		resp := MakeRequest(t, "DELETE", "/texts/"+created.ID, nil)
		resp.Body.Close()
	}()

	// Update the text
	updateReq := UpdateTextRequest{
		Content:       "Updated content",
		PageSlug:      "updated-page",
		LastUpdatedBy: "integration-test",
	}

	resp = MakeRequest(t, "PUT", "/texts/"+created.ID, updateReq)
	AssertStatusCode(t, resp, http.StatusOK)

	var updated TextResponse
	ParseJSONResponse(t, resp, &updated)

	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Updated content", updated.Content)
	assert.Equal(t, "updated-page", updated.PageSlug)
	assert.Equal(t, "integration-test", updated.LastUpdatedBy)
}

func TestTexts_Delete(t *testing.T) {
	slug := GenerateUniqueSlug("delete-test")

	// Create a text
	createReq := CreateTextRequest{
		Slug:    slug,
		Content: "This will be deleted",
	}

	resp := MakeRequest(t, "POST", "/texts", createReq)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created TextResponse
	ParseJSONResponse(t, resp, &created)

	// Delete the text
	resp = MakeRequest(t, "DELETE", "/texts/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusNoContent)
	resp.Body.Close()

	// Verify it's deleted
	resp = MakeRequest(t, "GET", "/texts/id/"+created.ID, nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestTexts_List(t *testing.T) {
	// Create multiple texts
	texts := []CreateTextRequest{
		{Slug: GenerateUniqueSlug("list-test-1"), Content: "Content 1", PageSlug: "test-page"},
		{Slug: GenerateUniqueSlug("list-test-2"), Content: "Content 2", PageSlug: "test-page"},
	}

	var createdIDs []string
	for _, text := range texts {
		resp := MakeRequest(t, "POST", "/texts", text)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var created TextResponse
		ParseJSONResponse(t, resp, &created)
		createdIDs = append(createdIDs, created.ID)
	}

	// Cleanup
	defer func() {
		for _, id := range createdIDs {
			resp := MakeRequest(t, "DELETE", "/texts/"+id, nil)
			resp.Body.Close()
		}
	}()

	// List all texts
	resp := MakeRequest(t, "GET", "/texts", nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var allTexts []TextResponse
	ParseJSONResponse(t, resp, &allTexts)

	assert.GreaterOrEqual(t, len(allTexts), 2, "Should have at least our 2 created texts")
}

func TestTexts_GetByPageSlug(t *testing.T) {
	pageSlug := GenerateUniqueSlug("test-page")

	// Create texts with the same page slug
	texts := []CreateTextRequest{
		{Slug: GenerateUniqueSlug("page-test-1"), Content: "Content 1", PageSlug: pageSlug},
		{Slug: GenerateUniqueSlug("page-test-2"), Content: "Content 2", PageSlug: pageSlug},
	}

	var createdIDs []string
	for _, text := range texts {
		resp := MakeRequest(t, "POST", "/texts", text)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var created TextResponse
		ParseJSONResponse(t, resp, &created)
		createdIDs = append(createdIDs, created.ID)
	}

	// Cleanup
	defer func() {
		for _, id := range createdIDs {
			resp := MakeRequest(t, "DELETE", "/texts/"+id, nil)
			resp.Body.Close()
		}
	}()

	// Get texts by page slug
	resp := MakeRequest(t, "GET", "/texts/page/slug/"+pageSlug, nil)
	AssertStatusCode(t, resp, http.StatusOK)

	var pageTexts []TextResponse
	ParseJSONResponse(t, resp, &pageTexts)

	assert.GreaterOrEqual(t, len(pageTexts), 2, "Should have at least our 2 texts for this page")

	// Verify all have the correct page slug
	for _, text := range pageTexts {
		if text.ID == createdIDs[0] || text.ID == createdIDs[1] {
			assert.Equal(t, pageSlug, text.PageSlug)
		}
	}
}

func TestTexts_NotFound(t *testing.T) {
	// Try to get non-existent text
	resp := MakeRequest(t, "GET", "/texts/non-existent-slug-12345", nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()

	// Try to get by non-existent ID
	resp = MakeRequest(t, "GET", "/texts/id/non-existent-id-12345", nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()

	// Try to update non-existent text
	updateReq := UpdateTextRequest{Content: "Updated"}
	resp = MakeRequest(t, "PUT", "/texts/non-existent-id-12345", updateReq)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()

	// Try to delete non-existent text
	resp = MakeRequest(t, "DELETE", "/texts/non-existent-id-12345", nil)
	AssertStatusCode(t, resp, http.StatusNotFound)
	resp.Body.Close()
}
