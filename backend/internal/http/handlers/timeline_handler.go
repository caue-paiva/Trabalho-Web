package handlers

import (
	"encoding/json"
	"net/http"

	"backend/internal/http/mapper"
	"backend/internal/platform/httputil"
	"backend/internal/service"
)

type TimelineHandler struct {
	textContentService service.TextContentService
}

func NewTimelineHandler(svc service.TextContentService) *TimelineHandler {
	return &TimelineHandler{textContentService: svc}
}

// ListTimelineEntries handles GET /api/v1/timelineentries
func (h *TimelineHandler) ListTimelineEntries(w http.ResponseWriter, r *http.Request) {
	entries, err := h.textContentService.ListTimelineEntries(r.Context())
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.TimelineEntriesToResponse(entries)
	httputil.JSON(w, response, http.StatusOK)
}

// GetTimelineEntryByID handles GET /api/v1/timelineentries/{id}
func (h *TimelineHandler) GetTimelineEntryByID(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "id")

	entry, err := h.textContentService.GetTimelineEntryByID(r.Context(), id)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.TimelineEntryToResponse(entry)
	httputil.JSON(w, response, http.StatusOK)
}

// CreateTimelineEntry handles POST /api/v1/timelineentries
func (h *TimelineHandler) CreateTimelineEntry(w http.ResponseWriter, r *http.Request) {
	var req mapper.CreateTimelineEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, err, http.StatusBadRequest)
		return
	}

	entity, err := mapper.ToTimelineEntryEntity(req)
	if err != nil {
		httputil.Error(w, err, http.StatusBadRequest)
		return
	}

	created, err := h.textContentService.CreateTimelineEntry(r.Context(), entity)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.TimelineEntryToResponse(created)
	httputil.JSON(w, response, http.StatusCreated)
}

// UpdateTimelineEntry handles PUT /api/v1/timelineentries/{id}
func (h *TimelineHandler) UpdateTimelineEntry(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "id")

	var req mapper.UpdateTimelineEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, err, http.StatusBadRequest)
		return
	}

	entity, err := mapper.ToTimelineEntryUpdateEntity(req)
	if err != nil {
		httputil.Error(w, err, http.StatusBadRequest)
		return
	}

	updated, err := h.textContentService.UpdateTimelineEntry(r.Context(), id, entity)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.TimelineEntryToResponse(updated)
	httputil.JSON(w, response, http.StatusOK)
}

// DeleteTimelineEntry handles DELETE /api/v1/timelineentries/{id}
func (h *TimelineHandler) DeleteTimelineEntry(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "id")

	if err := h.textContentService.DeleteTimelineEntry(r.Context(), id); err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.NoContent(w)
}
