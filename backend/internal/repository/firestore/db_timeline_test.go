package firestore

import (
	"context"
	"testing"
	"time"

	"backend/internal/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBRepository_CreateTimelineEntry(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name         string
		entry        entities.TimelineEntry
		expectError  bool
		validateFunc func(t *testing.T, created entities.TimelineEntry)
	}{
		{
			name: "create entry with all fields",
			entry: entities.TimelineEntry{
				Name:          "Grupy Sanca Foundation",
				Text:          "The Python community in São Carlos was officially founded",
				Location:      "São Carlos, SP",
				Date:          time.Date(2015, 6, 20, 0, 0, 0, 0, time.UTC),
				LastUpdatedBy: "admin",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			expectError: false,
			validateFunc: func(t *testing.T, created entities.TimelineEntry) {
				assert.NotEmpty(t, created.ID, "Should have an ID")
				assert.Equal(t, "Grupy Sanca Foundation", created.Name)
				assert.Equal(t, "The Python community in São Carlos was officially founded", created.Text)
				assert.Equal(t, "São Carlos, SP", created.Location)
				assert.Equal(t, "admin", created.LastUpdatedBy)
				assert.False(t, created.Date.IsZero(), "Date should be set")
			},
		},
		{
			name: "create entry without optional location",
			entry: entities.TimelineEntry{
				Name:      "Community Milestone",
				Text:      "Reached 100 active members",
				Date:      time.Date(2017, 8, 10, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
			validateFunc: func(t *testing.T, created entities.TimelineEntry) {
				assert.NotEmpty(t, created.ID, "Should have an ID")
				assert.Equal(t, "Community Milestone", created.Name)
				assert.Equal(t, "Reached 100 active members", created.Text)
				assert.Empty(t, created.Location, "Location should be empty")
				assert.False(t, created.Date.IsZero(), "Date should be set")
			},
		},
		{
			name: "create entry with minimal required fields",
			entry: entities.TimelineEntry{
				Name:      "Test Event",
				Text:      "Test event description",
				Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
			validateFunc: func(t *testing.T, created entities.TimelineEntry) {
				assert.NotEmpty(t, created.ID, "Should have an ID")
				assert.Equal(t, "Test Event", created.Name)
				assert.Equal(t, "Test event description", created.Text)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the entry
			created, err := db.CreateTimelineEntry(ctx, tt.entry)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err, "Failed to create timeline entry")

			// Cleanup
			defer func() {
				err := db.DeleteTimelineEntry(ctx, created.ID)
				assert.NoError(t, err, "Failed to cleanup created entry")
			}()

			// Run validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, created)
			}

			// Verify we can retrieve it
			retrieved, err := db.GetTimelineEntryByID(ctx, created.ID)
			require.NoError(t, err, "Failed to get created entry")
			assert.Equal(t, created.ID, retrieved.ID)
			assert.Equal(t, created.Name, retrieved.Name)
		})
	}
}

func TestDBRepository_GetTimelineEntryByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test entry for successful retrieval
	testEntry := entities.TimelineEntry{
		Name:      "First Python Workshop",
		Text:      "Our first educational workshop attracted 50+ participants",
		Location:  "UFSCar - São Carlos",
		Date:      time.Date(2016, 3, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := db.CreateTimelineEntry(ctx, testEntry)
	require.NoError(t, err, "Failed to create test entry")

	defer func() {
		db.DeleteTimelineEntry(ctx, created.ID)
	}()

	tests := []struct {
		name         string
		entryID      string
		expectError  bool
		errorMsg     string
		validateFunc func(t *testing.T, entry entities.TimelineEntry)
	}{
		{
			name:        "get existing entry by ID",
			entryID:     created.ID,
			expectError: false,
			validateFunc: func(t *testing.T, entry entities.TimelineEntry) {
				assert.Equal(t, created.ID, entry.ID)
				assert.Equal(t, "First Python Workshop", entry.Name)
				assert.Equal(t, "Our first educational workshop attracted 50+ participants", entry.Text)
				assert.Equal(t, "UFSCar - São Carlos", entry.Location)
			},
		},
		{
			name:        "get non-existent entry by ID",
			entryID:     "non-existent-entry-id-12345",
			expectError: true,
			errorMsg:    "not found",
		},
		{
			name:        "get with empty ID",
			entryID:     "",
			expectError: true,
			errorMsg:    "", // Don't check specific error message for invalid input
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := db.GetTimelineEntryByID(ctx, tt.entryID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				return
			}

			require.NoError(t, err, "Failed to get entry")
			if tt.validateFunc != nil {
				tt.validateFunc(t, entry)
			}
		})
	}
}

func TestDBRepository_UpdateTimelineEntry(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name         string
		setupEntry   entities.TimelineEntry
		updatePatch  entities.TimelineEntry
		expectError  bool
		errorMsg     string
		validateFunc func(t *testing.T, original, updated entities.TimelineEntry)
	}{
		{
			name: "full update of all fields",
			setupEntry: entities.TimelineEntry{
				Name:      "Original Event",
				Text:      "Original description",
				Location:  "Original Location",
				Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			updatePatch: entities.TimelineEntry{
				Name:          "Updated Event",
				Text:          "Updated description with more details",
				Location:      "Updated Location, SP",
				Date:          time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC),
				LastUpdatedBy: "test-user",
			},
			expectError: false,
			validateFunc: func(t *testing.T, original, updated entities.TimelineEntry) {
				assert.Equal(t, "Updated Event", updated.Name)
				assert.Equal(t, "Updated description with more details", updated.Text)
				assert.Equal(t, "Updated Location, SP", updated.Location)
				assert.Equal(t, "test-user", updated.LastUpdatedBy)
				assert.NotEqual(t, original.Date, updated.Date, "Date should be updated")
			},
		},
		{
			name: "partial update - only text and lastUpdatedBy",
			setupEntry: entities.TimelineEntry{
				Name:      "First Python Workshop",
				Text:      "Original workshop description",
				Location:  "UFSCar - São Carlos",
				Date:      time.Date(2016, 3, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			updatePatch: entities.TimelineEntry{
				Text:          "Updated: Our first educational workshop was a huge success with 50+ participants learning Python basics",
				LastUpdatedBy: "admin-user",
			},
			expectError: false,
			validateFunc: func(t *testing.T, original, updated entities.TimelineEntry) {
				assert.Equal(t, "Updated: Our first educational workshop was a huge success with 50+ participants learning Python basics", updated.Text)
				assert.Equal(t, "admin-user", updated.LastUpdatedBy)
				// Other fields should remain unchanged
				assert.Equal(t, original.Name, updated.Name)
				assert.Equal(t, original.Location, updated.Location)
			},
		},
		{
			name: "update name only",
			setupEntry: entities.TimelineEntry{
				Name:      "Old Name",
				Text:      "Some text",
				Location:  "Some location",
				Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			updatePatch: entities.TimelineEntry{
				Name: "Major Community Milestone",
			},
			expectError: false,
			validateFunc: func(t *testing.T, original, updated entities.TimelineEntry) {
				assert.Equal(t, "Major Community Milestone", updated.Name)
				assert.Equal(t, original.Text, updated.Text)
				assert.Equal(t, original.Location, updated.Location)
			},
		},
		{
			name: "add location to entry that had none",
			setupEntry: entities.TimelineEntry{
				Name:      "Community Milestone",
				Text:      "Reached 100 active members",
				Date:      time.Date(2017, 8, 10, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			updatePatch: entities.TimelineEntry{
				Location:      "São Carlos, Brazil",
				LastUpdatedBy: "admin",
			},
			expectError: false,
			validateFunc: func(t *testing.T, original, updated entities.TimelineEntry) {
				assert.Empty(t, original.Location, "Original should have no location")
				assert.Equal(t, "São Carlos, Brazil", updated.Location, "Updated should have location")
				assert.Equal(t, original.Name, updated.Name)
				assert.Equal(t, original.Text, updated.Text)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the initial entry
			created, err := db.CreateTimelineEntry(ctx, tt.setupEntry)
			require.NoError(t, err, "Failed to create setup entry")

			defer func() {
				db.DeleteTimelineEntry(ctx, created.ID)
			}()

			// Perform update
			updated, err := db.UpdateTimelineEntry(ctx, created.ID, tt.updatePatch)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				return
			}

			require.NoError(t, err, "Failed to update entry")

			// Run validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, created, updated)
			}

			// Verify by fetching again
			retrieved, err := db.GetTimelineEntryByID(ctx, created.ID)
			require.NoError(t, err, "Failed to get updated entry")
			assert.Equal(t, updated.Name, retrieved.Name)
			assert.Equal(t, updated.Text, retrieved.Text)
		})
	}
}

func TestDBRepository_UpdateTimelineEntry_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	patch := entities.TimelineEntry{
		Name: "Updated Name",
		Text: "Updated Text",
	}

	_, err := db.UpdateTimelineEntry(ctx, "non-existent-entry-id-12345", patch)
	assert.Error(t, err, "Should return error when updating non-existent entry")
	assert.Contains(t, err.Error(), "not found", "Error should mention 'not found'")
}

func TestDBRepository_ListTimelineEntries(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test entries with different dates
	testEntries := []entities.TimelineEntry{
		{
			Name:      "2015 - Grupy Sanca Foundation",
			Text:      "The Python community in São Carlos was officially founded",
			Location:  "São Carlos, SP",
			Date:      time.Date(2015, 6, 20, 0, 0, 0, 0, time.UTC),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "2024 - Python Conference",
			Text:      "Hosted regional Python conference with international speakers",
			Location:  "São Carlos Convention Center",
			Date:      time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "2017 - Community Milestone",
			Text:      "Reached 100 active members",
			Location:  "São Carlos, Brazil",
			Date:      time.Date(2017, 8, 10, 0, 0, 0, 0, time.UTC),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	var createdIDs []string
	for _, entry := range testEntries {
		created, err := db.CreateTimelineEntry(ctx, entry)
		require.NoError(t, err, "Failed to create entry")
		createdIDs = append(createdIDs, created.ID)
	}

	defer func() {
		for _, id := range createdIDs {
			db.DeleteTimelineEntry(ctx, id)
		}
	}()

	// List all entries
	entries, err := db.ListTimelineEntries(ctx)
	require.NoError(t, err, "Failed to list timeline entries")
	assert.GreaterOrEqual(t, len(entries), 3, "Should have at least 3 entries")

	// Find our test entries and verify chronological ordering
	var ourEntries []entities.TimelineEntry
	for _, entry := range entries {
		for _, id := range createdIDs {
			if entry.ID == id {
				ourEntries = append(ourEntries, entry)
			}
		}
	}

	require.Len(t, ourEntries, 3, "Should find all 3 created entries")

	// Verify chronological ordering (oldest first)
	assert.True(t, ourEntries[0].Date.Before(ourEntries[1].Date), "Entries should be in chronological order")
	assert.True(t, ourEntries[1].Date.Before(ourEntries[2].Date), "Entries should be in chronological order")

	// Verify specific order: 2015 -> 2017 -> 2024
	assert.Equal(t, "2015 - Grupy Sanca Foundation", ourEntries[0].Name)
	assert.Equal(t, "2017 - Community Milestone", ourEntries[1].Name)
	assert.Equal(t, "2024 - Python Conference", ourEntries[2].Name)
}

func TestDBRepository_DeleteTimelineEntry(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test entry
	newEntry := entities.TimelineEntry{
		Name:      "Test Event to Delete",
		Text:      "This entry will be deleted",
		Location:  "São Carlos, SP",
		Date:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := db.CreateTimelineEntry(ctx, newEntry)
	require.NoError(t, err, "Failed to create entry")

	// Delete the entry
	err = db.DeleteTimelineEntry(ctx, created.ID)
	require.NoError(t, err, "Failed to delete entry")

	// Verify it's deleted by trying to get it
	_, err = db.GetTimelineEntryByID(ctx, created.ID)
	assert.Error(t, err, "Should get error when fetching deleted entry")
	assert.Contains(t, err.Error(), "not found", "Error should mention 'not found'")
}

func TestDBRepository_ListTimelineEntries_EmptyResult(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// This test assumes we can query and get at least an empty array
	// Even if there are entries in the DB, this should not error
	entries, err := db.ListTimelineEntries(ctx)
	require.NoError(t, err, "Should not return error for listing")
	assert.NotNil(t, entries, "Should return a slice (even if empty)")
}
