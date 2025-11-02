package handlers

import (
	"encoding/json"
	"net/http"

	"backend/internal/http/mapper"
	"backend/internal/platform/httputil"
	"backend/internal/server"
)

type TextsHandler struct {
	server server.Server
}

func NewTextsHandler(srv server.Server) *TextsHandler {
	return &TextsHandler{server: srv}
}

// ListTexts handles GET /api/v1/texts
func (h *TextsHandler) ListTexts(w http.ResponseWriter, r *http.Request) {
	texts, err := h.server.ListAllTexts(r.Context())
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.TextsToResponse(texts)
	httputil.JSON(w, response, http.StatusOK)
}

// GetTextBySlug handles GET /api/v1/texts/{slug}
func (h *TextsHandler) GetTextBySlug(w http.ResponseWriter, r *http.Request) {
	slug := extractPathParam(r, "slug")

	text, err := h.server.GetTextBySlug(r.Context(), slug)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.TextToResponse(text)
	httputil.JSON(w, response, http.StatusOK)
}

// GetTextByID handles GET /api/v1/texts/id/{id}
func (h *TextsHandler) GetTextByID(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "id")

	text, err := h.server.GetTextByID(r.Context(), id)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.TextToResponse(text)
	httputil.JSON(w, response, http.StatusOK)
}

// GetTextsByPageID handles GET /api/v1/texts/page/{pageId}
func (h *TextsHandler) GetTextsByPageID(w http.ResponseWriter, r *http.Request) {
	pageID := extractPathParam(r, "pageId")

	texts, err := h.server.GetTextsByPageID(r.Context(), pageID)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.TextsToResponse(texts)
	httputil.JSON(w, response, http.StatusOK)
}

// GetTextsByPageSlug handles GET /api/v1/texts/page/slug/{pageSlug}
func (h *TextsHandler) GetTextsByPageSlug(w http.ResponseWriter, r *http.Request) {
	pageSlug := extractPathParam(r, "pageSlug")

	texts, err := h.server.GetTextsByPageSlug(r.Context(), pageSlug)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.TextsToResponse(texts)
	httputil.JSON(w, response, http.StatusOK)
}

// CreateText handles POST /api/v1/texts
func (h *TextsHandler) CreateText(w http.ResponseWriter, r *http.Request) {
	var req mapper.CreateTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, err, http.StatusBadRequest)
		return
	}

	entity := mapper.ToTextEntity(req)
	created, err := h.server.CreateText(r.Context(), entity)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.TextToResponse(created)
	httputil.JSON(w, response, http.StatusCreated)
}

// UpdateText handles PUT /api/v1/texts/{id}
func (h *TextsHandler) UpdateText(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "id")

	var req mapper.UpdateTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, err, http.StatusBadRequest)
		return
	}

	entity := mapper.ToTextUpdateEntity(req)
	updated, err := h.server.UpdateText(r.Context(), id, entity)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.TextToResponse(updated)
	httputil.JSON(w, response, http.StatusOK)
}

// DeleteText handles DELETE /api/v1/texts/{id}
func (h *TextsHandler) DeleteText(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "id")

	if err := h.server.DeleteText(r.Context(), id); err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.NoContent(w)
}
