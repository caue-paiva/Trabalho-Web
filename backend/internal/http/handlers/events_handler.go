package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/http/mapper"
	"backend/internal/platform/httputil"
	"backend/internal/server"
)

type EventsHandler struct {
	server server.Server
}

func NewEventsHandler(srv server.Server) *EventsHandler {
	return &EventsHandler{server: srv}
}

// GetEvents handles GET /api/v1/events?limit=N&orderBy=field&desc=true
func (h *EventsHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	limit := 10
	if limitStr := query.Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	orderBy := query.Get("orderBy")
	if orderBy == "" {
		orderBy = "startDate"
	}

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
