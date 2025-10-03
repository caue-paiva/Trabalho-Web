package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"backend/internal/entities"
)

const (
	grupyBaseURL  = "https://eventos.grupysanca.com.br/api/v1"
	jsonAPIAccept = "application/vnd.api+json"
)

// Filter represents a single filter condition for the API
// Example: {"name":"starts-at","op":"lt","val":"2025-10-03T21:00:00Z"}
type Filter struct {
	Name string `json:"name"` // Field name (e.g., "starts-at", "ends-at", "name")
	Op   string `json:"op"`   // Operator (e.g., "eq", "ne", "lt", "le", "gt", "ge", "like", "in")
	Val  string `json:"val"`  // Value to compare against
}

// QueryParams represents query parameters for the Grupy Sanca Events API
type QueryParams struct {
	Sort       string   // Sort field (e.g., "starts-at", "-starts-at" for descending)
	PageSize   int      // Number of results per page (max 100)
	PageNumber int      // Page number (1-based)
	Filters    []Filter // Array of filter conditions
}

// GrupyEventsAPI implements the gateway for Grupy Sanca Events API
type GrupyEventsAPI struct {
	client  *http.Client
	baseURL string
}

// NewGrupyEventsAPI creates a new Grupy Sanca Events API gateway
func NewGrupyEventsAPI() *GrupyEventsAPI {
	return &GrupyEventsAPI{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: grupyBaseURL,
	}
}

// JSON:API response structures

type jsonAPIResponse struct {
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
	Data []jsonAPIEventData `json:"data"`
}

type jsonAPIEventData struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Attributes jsonAPIEventAttrs `json:"attributes"`
}

type jsonAPIEventAttrs struct {
	Name              string  `json:"name"`
	Description       *string `json:"description"`
	StartsAt          string  `json:"starts-at"`
	EndsAt            string  `json:"ends-at"`
	Timezone          string  `json:"timezone"`
	LocationName      *string `json:"location-name"`
	LogoURL           *string `json:"logo-url"`
	ThumbnailImageURL *string `json:"thumbnail-image-url"`
	LargeImageURL     *string `json:"large-image-url"`
	OriginalImageURL  *string `json:"original-image-url"`
	IconImageURL      *string `json:"icon-image-url"`
	Identifier        string  `json:"identifier"`
	Privacy           string  `json:"privacy"`
	State             string  `json:"state"`
	CreatedAt         string  `json:"created-at"`
}

// GetEvents fetches events from Grupy Sanca API using QueryParams
func (g *GrupyEventsAPI) GetEvents(ctx context.Context, params QueryParams) ([]entities.Event, error) {
	// Build URL with query parameters
	apiURL, err := g.buildEventsURLFromParams(params)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set JSON:API headers
	req.Header.Set("Accept", jsonAPIAccept)

	// Execute request
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse JSON:API response
	var apiResp jsonAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Map to entities
	events := make([]entities.Event, 0, len(apiResp.Data))
	for _, data := range apiResp.Data {
		event, err := g.mapToEntity(data)
		if err != nil {
			// Log error but continue processing other events
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

// buildEventsURLFromParams constructs the API URL with QueryParams
func (g *GrupyEventsAPI) buildEventsURLFromParams(params QueryParams) (string, error) {
	endpoint := g.baseURL + "/events"
	queryParams := url.Values{}

	// Add filters as JSON array
	if len(params.Filters) > 0 {
		filtersJSON, err := json.Marshal(params.Filters)
		if err != nil {
			return "", fmt.Errorf("failed to marshal filters: %w", err)
		}
		queryParams.Set("filter", string(filtersJSON))
	}

	// Add sort parameter
	if params.Sort != "" {
		queryParams.Set("sort", params.Sort)
	}

	// Add pagination
	if params.PageSize > 0 {
		if params.PageSize > 100 {
			params.PageSize = 100 // API max limit
		}
		queryParams.Set("page[size]", strconv.Itoa(params.PageSize))
	}

	if params.PageNumber > 0 {
		queryParams.Set("page[number]", strconv.Itoa(params.PageNumber))
	}

	// Build final URL
	if len(queryParams) > 0 {
		endpoint = endpoint + "?" + queryParams.Encode()
	}

	return endpoint, nil
}

// mapSortField maps our sort field names to API field names
func (g *GrupyEventsAPI) mapSortField(orderBy string) string {
	switch orderBy {
	case "startDate":
		return "starts-at"
	case "name":
		return "name"
	case "created":
		return "created-at"
	default:
		return "starts-at"
	}
}

// mapToEntity maps JSON:API event data to our Event entity
func (g *GrupyEventsAPI) mapToEntity(data jsonAPIEventData) (entities.Event, error) {
	attrs := data.Attributes

	// Parse timestamps
	startsAt, err := time.Parse(time.RFC3339, attrs.StartsAt)
	if err != nil {
		return entities.Event{}, fmt.Errorf("invalid starts-at: %w", err)
	}

	endsAt, err := time.Parse(time.RFC3339, attrs.EndsAt)
	if err != nil {
		return entities.Event{}, fmt.Errorf("invalid ends-at: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, attrs.CreatedAt)
	if err != nil {
		// CreatedAt is optional, use zero value if parsing fails
		createdAt = time.Time{}
	}

	// Build event entity
	event := entities.Event{
		ID:         data.ID,
		Identifier: attrs.Identifier,
		Name:       attrs.Name,
		StartsAt:   startsAt,
		EndsAt:     endsAt,
		Timezone:   attrs.Timezone,
		Privacy:    attrs.Privacy,
		State:      attrs.State,
		CreatedAt:  createdAt,
	}

	// Handle optional fields
	if attrs.Description != nil {
		event.Description = *attrs.Description
	}
	if attrs.LocationName != nil {
		event.LocationName = *attrs.LocationName
	}
	if attrs.LogoURL != nil {
		event.LogoURL = *attrs.LogoURL
	}
	if attrs.ThumbnailImageURL != nil {
		event.ThumbnailImageURL = *attrs.ThumbnailImageURL
	}
	if attrs.LargeImageURL != nil {
		event.LargeImageURL = *attrs.LargeImageURL
	}
	if attrs.OriginalImageURL != nil {
		event.OriginalImageURL = *attrs.OriginalImageURL
	}
	if attrs.IconImageURL != nil {
		event.IconImageURL = *attrs.IconImageURL
	}

	return event, nil
}
