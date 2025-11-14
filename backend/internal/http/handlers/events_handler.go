package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/http/mapper"
	"backend/internal/platform/httputil"
)

// GetEvents handles GET /api/v1/events?limit=N&orderBy=starts-at&desc=true
// Follows the same logic as the Grupy API query and filter field names: starts-at, ends-at, name, created-at, etc.
func (h *BaseHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	limit := 10
	if limitStr := query.Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	orderBy := query.Get("orderBy") // Pass through, default handled at client level

	desc := false
	if descStr := query.Get("desc"); descStr == "true" {
		desc = true
	}

	// Call service
	events, err := h.server.GetEvents(r.Context(), limit, orderBy, desc)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.EventsToResponse(events)
	httputil.JSON(w, response, http.StatusOK)
}
