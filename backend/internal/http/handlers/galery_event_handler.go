package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend/internal/http/mapper"
	"backend/internal/platform/httputil"
)

// CreateGaleryEvent handles POST /api/v1/galery_events
func (h *BaseHandler) CreateGaleryEvent(w http.ResponseWriter, r *http.Request) {
	var req mapper.CreateGaleryEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, err, http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" {
		httputil.Error(w, fmt.Errorf("name is required"), http.StatusBadRequest)
		return
	}
	if req.Location == "" {
		httputil.Error(w, fmt.Errorf("location is required"), http.StatusBadRequest)
		return
	}
	if req.Date.IsZero() {
		httputil.Error(w, fmt.Errorf("date is required"), http.StatusBadRequest)
		return
	}
	if len(req.ImagesBase64) == 0 {
		httputil.Error(w, fmt.Errorf("at least one image is required"), http.StatusBadRequest)
		return
	}

	// Create galery event (uploads images and saves to DB)
	created, err := h.server.CreateGaleryEvent(
		r.Context(),
		req.Name,
		req.Location,
		req.Date,
		req.ImagesBase64,
	)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.GaleryEventToResponse(created)
	httputil.JSON(w, response, http.StatusCreated)
}

// GetGaleryEventByID handles GET /api/v1/galery_events/{id}
func (h *BaseHandler) GetGaleryEventByID(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "id")

	event, err := h.server.GetGaleryEventByID(r.Context(), id)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.GaleryEventToResponse(event)
	httputil.JSON(w, response, http.StatusOK)
}

// ListGaleryEvents handles GET /api/v1/galery_events
func (h *BaseHandler) ListGaleryEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.server.ListGaleryEvents(r.Context())
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.GaleryEventsToResponse(events)
	httputil.JSON(w, response, http.StatusOK)
}

// DeleteGaleryEvent handles DELETE /api/v1/galery_events/{id}
// Note: This deletes only the database record, not the associated images
func (h *BaseHandler) DeleteGaleryEvent(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "id")

	if err := h.server.DeleteGaleryEvent(r.Context(), id); err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
