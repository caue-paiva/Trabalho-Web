package firestore

import (
	"context"
	"os"
	"testing"
	"time"

	"backend/configs"
	"backend/internal/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getTestCollections loads collection names from config
func getTestCollections(t *testing.T) CollectionNames {
	os.Unsetenv("RUNTIME_ENV")

	config, err := configs.NewConfigService()
	require.NoError(t, err, "Failed to load config")

	type Collections struct {
		Texts     string `yaml:"texts"`
		Images    string `yaml:"images"`
		Timelines string `yaml:"timelines"`
	}

	var cols Collections
	err = config.UnmarshalKey("collections", &cols)
	require.NoError(t, err, "Failed to unmarshal collections")

	return CollectionNames{
		Texts:           cols.Texts,
		Images:          cols.Images,
		TimelineEntries: cols.Timelines,
	}
}

// getTestFirestoreConfig loads complete Firestore config from development.yaml
func getTestFirestoreConfig(t *testing.T) FirestoreConfig {
	os.Unsetenv("RUNTIME_ENV")

	config, err := configs.NewConfigService()
	require.NoError(t, err, "Failed to load config")

	type FirebaseConfig struct {
		ProjectID       string `yaml:"project_id"`
		CredentialsPath string `yaml:"credentials_path"`
	}

	var fbConfig FirebaseConfig
	err = config.UnmarshalKey("firebase", &fbConfig)
	require.NoError(t, err, "Failed to unmarshal firebase config")

	// Get credentials JSON bytes
	credentialsJSON, err := config.GetCredentialsJSON(fbConfig.CredentialsPath)
	require.NoError(t, err, "Failed to get credentials JSON")

	collections := getTestCollections(t)

	return FirestoreConfig{
		ProjectID:       fbConfig.ProjectID,
		CredentialsJSON: credentialsJSON,
		Collections:     collections,
	}
}

// setupTestDB creates a test DB repository
func setupTestDB(t *testing.T) (*DBRepository, func()) {
	ctx := context.Background()
	config := getTestFirestoreConfig(t)

	client, err := NewFirestoreClient(ctx, config)
	require.NoError(t, err, "Failed to create Firestore client")

	db := NewDBRepository(client, config.Collections)

	cleanup := func() {
		client.Close()
	}

	return db, cleanup
}

func TestDBRepository_CreateAndGetText(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test text entry
	newText := entities.Text{
		Slug:      "test-create-get",
		Content:   "Test content for create and get",
		PageSlug:  "test-page",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create the text
	created, err := db.CreateText(ctx, newText)
	require.NoError(t, err, "Failed to create text")
	assert.NotEmpty(t, created.ID, "Created text should have an ID")
	assert.Equal(t, newText.Slug, created.Slug)
	assert.Equal(t, newText.Content, created.Content)

	// Cleanup: delete the created text
	defer func() {
		err := db.DeleteText(ctx, created.ID)
		assert.NoError(t, err, "Failed to cleanup created text")
	}()

	// Get by slug
	retrieved, err := db.GetTextBySlug(ctx, created.Slug)
	require.NoError(t, err, "Failed to get text by slug")
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Slug, retrieved.Slug)
	assert.Equal(t, created.Content, retrieved.Content)

	// Get by ID
	retrievedByID, err := db.GetTextByID(ctx, created.ID)
	require.NoError(t, err, "Failed to get text by ID")
	assert.Equal(t, created.ID, retrievedByID.ID)
	assert.Equal(t, created.Slug, retrievedByID.Slug)
}

func TestDBRepository_GetTextBySlug(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple texts with unique slugs
	testTexts := []entities.Text{
		{
			Slug:      "about-us-page",
			Content:   "Information about our organization",
			PageSlug:  "about",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Slug:      "contact-info",
			Content:   "How to reach us",
			PageSlug:  "contact",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Slug:      "mission-statement",
			Content:   "Our mission and values",
			PageSlug:  "about",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	var createdIDs []string
	for _, text := range testTexts {
		created, err := db.CreateText(ctx, text)
		require.NoError(t, err, "Failed to create text")
		createdIDs = append(createdIDs, created.ID)
	}

	// Cleanup all created texts
	defer func() {
		for _, id := range createdIDs {
			err := db.DeleteText(ctx, id)
			assert.NoError(t, err, "Failed to cleanup text")
		}
	}()

	// Test: Query by each slug and verify correct text is returned
	for i, expectedText := range testTexts {
		retrieved, err := db.GetTextBySlug(ctx, expectedText.Slug)
		require.NoError(t, err, "Failed to get text by slug: %s", expectedText.Slug)

		// Verify it's the correct text
		assert.Equal(t, createdIDs[i], retrieved.ID, "Should retrieve text with correct ID")
		assert.Equal(t, expectedText.Slug, retrieved.Slug, "Should retrieve text with correct slug")
		assert.Equal(t, expectedText.Content, retrieved.Content, "Should retrieve text with correct content")
		assert.Equal(t, expectedText.PageSlug, retrieved.PageSlug, "Should retrieve text with correct page slug")
	}

	// Test: Query by non-existent slug should return error
	_, err := db.GetTextBySlug(ctx, "non-existent-slug-12345")
	assert.Error(t, err, "Should return error for non-existent slug")
}

func TestDBRepository_UpdateText(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test text entry
	newText := entities.Text{
		Slug:      "test-update",
		Content:   "Original content",
		PageSlug:  "test-page",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := db.CreateText(ctx, newText)
	require.NoError(t, err, "Failed to create text")

	// Cleanup
	defer func() {
		err := db.DeleteText(ctx, created.ID)
		assert.NoError(t, err, "Failed to cleanup created text")
	}()

	// Update the text
	updated := created
	updated.Content = "Updated content"
	updated.PageSlug = "updated-page"

	result, err := db.UpdateText(ctx, created.ID, updated)
	require.NoError(t, err, "Failed to update text")
	assert.Equal(t, "Updated content", result.Content)
	assert.Equal(t, "updated-page", result.PageSlug)

	// Verify update by fetching again
	retrieved, err := db.GetTextByID(ctx, created.ID)
	require.NoError(t, err, "Failed to get updated text")
	assert.Equal(t, "Updated content", retrieved.Content)
	assert.Equal(t, "updated-page", retrieved.PageSlug)
}

func TestDBRepository_ListTexts(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple test texts
	texts := []entities.Text{
		{
			Slug:      "test-list-1",
			Content:   "Content 1",
			PageSlug:  "page-1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Slug:      "test-list-2",
			Content:   "Content 2",
			PageSlug:  "page-2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	var createdIDs []string
	for _, text := range texts {
		created, err := db.CreateText(ctx, text)
		require.NoError(t, err, "Failed to create text")
		createdIDs = append(createdIDs, created.ID)
	}

	// Cleanup all created texts
	defer func() {
		for _, id := range createdIDs {
			err := db.DeleteText(ctx, id)
			assert.NoError(t, err, "Failed to cleanup text")
		}
	}()

	// List all texts
	allTexts, err := db.ListAllTexts(ctx)
	require.NoError(t, err, "Failed to list texts")
	assert.GreaterOrEqual(t, len(allTexts), 2, "Should have at least the 2 texts we created")

	// Verify our texts are in the list
	foundCount := 0
	for _, text := range allTexts {
		for _, id := range createdIDs {
			if text.ID == id {
				foundCount++
			}
		}
	}
	assert.Equal(t, len(texts), foundCount, "Should find all created texts in the list")
}

func TestDBRepository_DeleteText(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test text
	newText := entities.Text{
		Slug:      "test-delete",
		Content:   "Content to be deleted",
		PageSlug:  "test-page",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := db.CreateText(ctx, newText)
	require.NoError(t, err, "Failed to create text")

	// Delete the text
	err = db.DeleteText(ctx, created.ID)
	require.NoError(t, err, "Failed to delete text")

	// Verify it's deleted by trying to get it
	_, err = db.GetTextByID(ctx, created.ID)
	assert.Error(t, err, "Should get error when fetching deleted text")
}

func TestDBRepository_ListTextsByPageSlug(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	pageSlug := "test-page-filter"

	// Create texts with specific page slug
	texts := []entities.Text{
		{
			Slug:      "test-page-1",
			Content:   "Content 1",
			PageSlug:  pageSlug,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Slug:      "test-page-2",
			Content:   "Content 2",
			PageSlug:  pageSlug,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Slug:      "test-other-page",
			Content:   "Other content",
			PageSlug:  "other-page",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	var createdIDs []string
	for _, text := range texts {
		created, err := db.CreateText(ctx, text)
		require.NoError(t, err, "Failed to create text")
		createdIDs = append(createdIDs, created.ID)
	}

	// Cleanup
	defer func() {
		for _, id := range createdIDs {
			err := db.DeleteText(ctx, id)
			assert.NoError(t, err, "Failed to cleanup text")
		}
	}()

	// List texts by page slug
	pageTexts, err := db.ListTextsByPageSlug(ctx, pageSlug)
	require.NoError(t, err, "Failed to list texts by page slug")
	assert.GreaterOrEqual(t, len(pageTexts), 2, "Should have at least 2 texts for this page")

	// Verify all returned texts have the correct page slug
	for _, text := range pageTexts {
		// Only check texts we created
		for _, id := range createdIDs[:2] {
			if text.ID == id {
				assert.Equal(t, pageSlug, text.PageSlug, "Text should have correct page slug")
			}
		}
	}
}
