package handlers

import (
	"encoding/json"
	"net/http"

	"backend/internal/http/mapper"
	"backend/internal/platform/httputil"
	"backend/internal/service"
)

type ImagesHandler struct {
	contentService service.ContentService
}

func NewImagesHandler(svc service.ContentService) *ImagesHandler {
	return &ImagesHandler{contentService: svc}
}

// GetImageByID handles GET /api/v1/images/{id}
func (h *ImagesHandler) GetImageByID(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "id")

	img, err := h.contentService.GetImageByID(r.Context(), id)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.ImageToResponse(img)
	httputil.JSON(w, response, http.StatusOK)
}

// GetImagesByGallerySlug handles GET /api/v1/images/gallery/{slug}
func (h *ImagesHandler) GetImagesByGallerySlug(w http.ResponseWriter, r *http.Request) {
	slug := extractPathParam(r, "slug")

	images, err := h.contentService.GetImagesByGallerySlug(r.Context(), slug)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.ImagesToResponse(images)
	httputil.JSON(w, response, http.StatusOK)
}

// CreateImage handles POST /api/v1/images
func (h *ImagesHandler) CreateImage(w http.ResponseWriter, r *http.Request) {
	var req mapper.CreateImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, err, http.StatusBadRequest)
		return
	}

	meta, data, err := mapper.ToImageEntity(req)
	if err != nil {
		httputil.Error(w, err, http.StatusBadRequest)
		return
	}

	created, err := h.contentService.UploadImage(r.Context(), meta, data)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.ImageToResponse(created)
	httputil.JSON(w, response, http.StatusCreated)
}

// UpdateImage handles PUT /api/v1/images/{id}
func (h *ImagesHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "id")

	var req mapper.UpdateImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, err, http.StatusBadRequest)
		return
	}

	meta, data, err := mapper.ToImageUpdateEntity(req)
	if err != nil {
		httputil.Error(w, err, http.StatusBadRequest)
		return
	}

	updated, err := h.contentService.UpdateImage(r.Context(), id, meta, data)
	if err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	response := mapper.ImageToResponse(updated)
	httputil.JSON(w, response, http.StatusOK)
}

// DeleteImage handles DELETE /api/v1/images/{id}
func (h *ImagesHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	id := extractPathParam(r, "id")

	if err := h.contentService.DeleteImage(r.Context(), id); err != nil {
		httputil.ErrorFromDomain(w, err)
		return
	}

	httputil.NoContent(w)
}
