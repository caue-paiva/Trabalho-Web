package clients

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/entities"
)

// TestEventsClient_GetEvents_ContractValidation validates the contract with the real Grupy API
// This test calls the actual external API and validates:
// 1. Non-empty response is returned
// 2. All events can be parsed to our internal structs
// 3. Required fields are present and valid
func TestEventsClient_GetEvents_ContractValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Arrange
	client := NewEventsClient()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Act
	events, err := client.GetEvents(ctx, 10, "starts-at", false)

	// Assert - Critical checks that should stop the test
	require.NoError(t, err, "API call should not fail")
	require.NotEmpty(t, events, "API should return at least one event - empty response indicates contract break")

	// Validate structure of all events
	for i, event := range events {
		t.Run(event.Name, func(t *testing.T) {
			// Required fields validation
			assert.NotEmpty(t, event.ID, "Event[%d]: ID must not be empty", i)
			assert.NotEmpty(t, event.Name, "Event[%d]: Name must not be empty", i)
			assert.NotEmpty(t, event.Identifier, "Event[%d]: Identifier must not be empty", i)
			assert.False(t, event.StartsAt.IsZero(), "Event[%d]: StartsAt must be a valid timestamp", i)
			assert.False(t, event.EndsAt.IsZero(), "Event[%d]: EndsAt must be a valid timestamp", i)
			assert.NotEmpty(t, event.Timezone, "Event[%d]: Timezone must not be empty", i)

			// Validate StartsAt is before EndsAt
			assert.True(t, event.StartsAt.Before(event.EndsAt) || event.StartsAt.Equal(event.EndsAt),
				"Event[%d]: StartsAt (%v) should be before or equal to EndsAt (%v)",
				i, event.StartsAt, event.EndsAt)

			// Optional fields - just log presence
			if event.Description != "" {
				t.Logf("  ✓ Event[%d] has description (%d chars)", i, len(event.Description))
			}
			if event.LocationName != "" {
				t.Logf("  ✓ Event[%d] has location: %s", i, event.LocationName)
			}
			if event.LogoURL != "" {
				t.Logf("  ✓ Event[%d] has logo URL", i)
			}
			if !event.CreatedAt.IsZero() {
				t.Logf("  ✓ Event[%d] has createdAt: %v", i, event.CreatedAt)
			}
		})
	}

	t.Logf("✓ Successfully validated %d events from Grupy API", len(events))
}

// TestEventsClient_GetEvents_DifferentQueries tests various query parameter combinations
func TestEventsClient_GetEvents_DifferentQueries(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewEventsClient()

	tests := []struct {
		name    string
		limit   int
		orderBy string
		desc    bool
		wantErr bool
	}{
		{
			name:    "default query with limit 10",
			limit:   10,
			orderBy: "starts-at",
			desc:    false,
			wantErr: false,
		},
		{
			name:    "large limit (50)",
			limit:   50,
			orderBy: "starts-at",
			desc:    false,
			wantErr: false,
		},
		{
			name:    "descending order by start date",
			limit:   10,
			orderBy: "starts-at",
			desc:    true,
			wantErr: false,
		},
		{
			name:    "sort by name ascending",
			limit:   10,
			orderBy: "name",
			desc:    false,
			wantErr: false,
		},
		{
			name:    "sort by created date descending",
			limit:   10,
			orderBy: "created-at",
			desc:    true,
			wantErr: false,
		},
		{
			name:    "maximum limit (100)",
			limit:   100,
			orderBy: "starts-at",
			desc:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			events, err := client.GetEvents(ctx, tt.limit, tt.orderBy, tt.desc)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, events, "Should return at least one event")
			assert.LessOrEqual(t, len(events), tt.limit, "Should not exceed requested limit")

			// Validate sort order if we have multiple events
			if len(events) > 1 && tt.orderBy == "starts-at" {
				for i := 0; i < len(events)-1; i++ {
					if tt.desc {
						assert.True(t,
							events[i].StartsAt.After(events[i+1].StartsAt) || events[i].StartsAt.Equal(events[i+1].StartsAt),
							"Events should be sorted by starts-at descending")
					} else {
						assert.True(t,
							events[i].StartsAt.Before(events[i+1].StartsAt) || events[i].StartsAt.Equal(events[i+1].StartsAt),
							"Events should be sorted by starts-at ascending")
					}
				}
			}

			t.Logf("%s: Got %d events", tt.name, len(events))
		})
	}
}

// TestEventsClient_MapToEntity tests the internal mapping logic
func TestEventsClient_MapToEntity(t *testing.T) {
	client := &eventsClient{
		httpClient: nil, // not needed for this test
		baseURL:    grupyBaseURL,
	}

	tests := []struct {
		name    string
		input   jsonAPIEventData
		want    entities.Event
		wantErr bool
	}{
		{
			name: "complete event with all fields",
			input: jsonAPIEventData{
				Type: "event",
				ID:   "123",
				Attributes: jsonAPIEventAttrs{
					Name:              "Test Event",
					Description:       stringPtr("Test Description"),
					StartsAt:          "2025-12-01T10:00:00Z",
					EndsAt:            "2025-12-01T12:00:00Z",
					Timezone:          "UTC",
					LocationName:      stringPtr("Test Location"),
					LogoURL:           stringPtr("https://example.com/logo.png"),
					ThumbnailImageURL: stringPtr("https://example.com/thumb.png"),
					LargeImageURL:     stringPtr("https://example.com/large.png"),
					OriginalImageURL:  stringPtr("https://example.com/original.png"),
					IconImageURL:      stringPtr("https://example.com/icon.png"),
					Identifier:        "test-123",
					Privacy:           "public",
					State:             "published",
					CreatedAt:         "2025-11-01T08:00:00Z",
				},
			},
			want: entities.Event{
				ID:                "123",
				Identifier:        "test-123",
				Name:              "Test Event",
				Description:       "Test Description",
				StartsAt:          parseTime(t, "2025-12-01T10:00:00Z"),
				EndsAt:            parseTime(t, "2025-12-01T12:00:00Z"),
				Timezone:          "UTC",
				LocationName:      "Test Location",
				LogoURL:           "https://example.com/logo.png",
				ThumbnailImageURL: "https://example.com/thumb.png",
				LargeImageURL:     "https://example.com/large.png",
				OriginalImageURL:  "https://example.com/original.png",
				IconImageURL:      "https://example.com/icon.png",
				Privacy:           "public",
				State:             "published",
				CreatedAt:         parseTime(t, "2025-11-01T08:00:00Z"),
			},
			wantErr: false,
		},
		{
			name: "minimal event with only required fields",
			input: jsonAPIEventData{
				Type: "event",
				ID:   "456",
				Attributes: jsonAPIEventAttrs{
					Name:       "Minimal Event",
					StartsAt:   "2025-12-15T14:00:00Z",
					EndsAt:     "2025-12-15T16:00:00Z",
					Timezone:   "America/Sao_Paulo",
					Identifier: "minimal-456",
					Privacy:    "public",
					State:      "draft",
					CreatedAt:  "2025-11-15T10:00:00Z",
				},
			},
			want: entities.Event{
				ID:         "456",
				Identifier: "minimal-456",
				Name:       "Minimal Event",
				StartsAt:   parseTime(t, "2025-12-15T14:00:00Z"),
				EndsAt:     parseTime(t, "2025-12-15T16:00:00Z"),
				Timezone:   "America/Sao_Paulo",
				Privacy:    "public",
				State:      "draft",
				CreatedAt:  parseTime(t, "2025-11-15T10:00:00Z"),
			},
			wantErr: false,
		},
		{
			name: "invalid starts-at timestamp",
			input: jsonAPIEventData{
				Type: "event",
				ID:   "789",
				Attributes: jsonAPIEventAttrs{
					Name:       "Invalid Event",
					StartsAt:   "invalid-date",
					EndsAt:     "2025-12-01T12:00:00Z",
					Timezone:   "UTC",
					Identifier: "invalid-789",
					Privacy:    "public",
					State:      "draft",
					CreatedAt:  "2025-11-01T08:00:00Z",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid ends-at timestamp",
			input: jsonAPIEventData{
				Type: "event",
				ID:   "790",
				Attributes: jsonAPIEventAttrs{
					Name:       "Invalid Event",
					StartsAt:   "2025-12-01T10:00:00Z",
					EndsAt:     "not-a-date",
					Timezone:   "UTC",
					Identifier: "invalid-790",
					Privacy:    "public",
					State:      "draft",
					CreatedAt:  "2025-11-01T08:00:00Z",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.mapToEntity(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Identifier, got.Identifier)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.True(t, tt.want.StartsAt.Equal(got.StartsAt), "StartsAt mismatch")
			assert.True(t, tt.want.EndsAt.Equal(got.EndsAt), "EndsAt mismatch")
			assert.Equal(t, tt.want.Timezone, got.Timezone)
			assert.Equal(t, tt.want.LocationName, got.LocationName)
			assert.Equal(t, tt.want.LogoURL, got.LogoURL)
			assert.Equal(t, tt.want.Privacy, got.Privacy)
			assert.Equal(t, tt.want.State, got.State)
		})
	}
}

// TestEventsClient_BuildSortParam tests the sort parameter building logic
func TestEventsClient_BuildSortParam(t *testing.T) {
	client := &eventsClient{}

	tests := []struct {
		name    string
		orderBy string
		desc    bool
		want    string
	}{
		{"starts-at ascending", "starts-at", false, "starts-at"},
		{"starts-at descending", "starts-at", true, "-starts-at"},
		{"name ascending", "name", false, "name"},
		{"name descending", "name", true, "-name"},
		{"created-at ascending", "created-at", false, "created-at"},
		{"created-at descending", "created-at", true, "-created-at"},
		{"ends-at ascending", "ends-at", false, "ends-at"},
		{"ends-at descending", "ends-at", true, "-ends-at"},
		{"empty field defaults to starts-at", "", false, "starts-at"},
		{"empty field descending", "", true, "-starts-at"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := client.buildSortParam(tt.orderBy, tt.desc)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}

func parseTime(t *testing.T, timeStr string) time.Time {
	parsed, err := time.Parse(time.RFC3339, timeStr)
	require.NoError(t, err)
	return parsed
}
