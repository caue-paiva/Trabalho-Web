package http

import (
	"net/http"

	"backend/internal/http/handlers"
	"backend/internal/platform/middleware"
	"backend/internal/service"
)

// NewRouter creates and configures the HTTP router
func NewRouter(server service.Server) http.Handler {
	mux := http.NewServeMux()

	// Create handlers
	textsHandler := handlers.NewTextsHandler(server)
	imagesHandler := handlers.NewImagesHandler(server)
	timelineHandler := handlers.NewTimelineHandler(server)
	eventsHandler := handlers.NewEventsHandler(server)

	// Register routes using Go 1.22+ pattern matching

	// Texts routes
	mux.HandleFunc("GET /api/v1/texts", textsHandler.ListTexts)
	mux.HandleFunc("GET /api/v1/texts/{slug}", textsHandler.GetTextBySlug)
	mux.HandleFunc("GET /api/v1/texts/id/{id}", textsHandler.GetTextByID)
	mux.HandleFunc("GET /api/v1/texts/page/{pageId}", textsHandler.GetTextsByPageID)
	mux.HandleFunc("GET /api/v1/texts/page/slug/{pageSlug}", textsHandler.GetTextsByPageSlug)
	mux.HandleFunc("POST /api/v1/texts", textsHandler.CreateText)
	mux.HandleFunc("PUT /api/v1/texts/{id}", textsHandler.UpdateText)
	mux.HandleFunc("DELETE /api/v1/texts/{id}", textsHandler.DeleteText)

	// Images routes
	mux.HandleFunc("GET /api/v1/images/{id}", imagesHandler.GetImageByID)
	mux.HandleFunc("GET /api/v1/images/gallery/{slug}", imagesHandler.GetImagesByGallerySlug)
	mux.HandleFunc("POST /api/v1/images", imagesHandler.CreateImage)
	mux.HandleFunc("PUT /api/v1/images/{id}", imagesHandler.UpdateImage)
	mux.HandleFunc("DELETE /api/v1/images/{id}", imagesHandler.DeleteImage)

	// Timeline routes
	mux.HandleFunc("GET /api/v1/timelineentries", timelineHandler.ListTimelineEntries)
	mux.HandleFunc("GET /api/v1/timelineentries/{id}", timelineHandler.GetTimelineEntryByID)
	mux.HandleFunc("POST /api/v1/timelineentries", timelineHandler.CreateTimelineEntry)
	mux.HandleFunc("PUT /api/v1/timelineentries/{id}", timelineHandler.UpdateTimelineEntry)
	mux.HandleFunc("DELETE /api/v1/timelineentries/{id}", timelineHandler.DeleteTimelineEntry)

	// Events routes
	mux.HandleFunc("GET /api/v1/events", eventsHandler.GetEvents)

	// Apply middleware (outermost to innermost)
	var handler http.Handler = mux
	handler = middleware.Recovery(handler)
	handler = middleware.CORS(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RequestID(handler)

	return handler
}
