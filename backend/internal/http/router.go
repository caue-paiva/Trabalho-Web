package http

import (
	"context"
	"log"
	"net/http"

	"backend/internal/http/handlers"
	"backend/internal/platform/auth"
	"backend/internal/platform/middleware"
	"backend/internal/server"
)

type RouterOptions struct {
	AuthConfig auth.AuthConfig
	Logger     *log.Logger
}

// NewRouter creates and configures the HTTP router
func NewRouter(ctx context.Context, srv server.Server, opts RouterOptions) http.Handler {
	mux := http.NewServeMux()

	// Create handlers
	textsHandler := handlers.NewBaseHandler(srv)
	imagesHandler := handlers.NewBaseHandler(srv)
	timelineHandler := handlers.NewBaseHandler(srv)
	eventsHandler := handlers.NewBaseHandler(srv)
	galeryEventHandler := handlers.NewBaseHandler(srv)
	authHandler := handlers.NewBaseHandler(srv)

	// Register routes using Go 1.22+ pattern matching

	// Texts routes
	mux.HandleFunc("GET /api/v1/texts", textsHandler.ListTexts)
	mux.HandleFunc("GET /api/v1/texts/{slug}", textsHandler.GetTextBySlug)
	mux.HandleFunc("GET /api/v1/texts/id/{id}", textsHandler.GetTextByID)
	mux.HandleFunc("GET /api/v1/texts/page/{pageId}", textsHandler.GetTextsByPageID)
	mux.HandleFunc("GET /api/v1/texts/page/slug/{pageSlug}", textsHandler.GetTextsByPageSlug)

	// Add auth middleware to non-get functions
	mux.HandleFunc("POST /api/v1/texts",
		middleware.NewAuthMiddlewareFunc(textsHandler.CreateText, opts.AuthConfig, opts.Logger),
	)
	mux.HandleFunc("PUT /api/v1/texts/{id}",
		middleware.NewAuthMiddlewareFunc(textsHandler.UpdateText, opts.AuthConfig, opts.Logger),
	)
	mux.HandleFunc("DELETE /api/v1/texts/{id}",
		middleware.NewAuthMiddlewareFunc(textsHandler.DeleteText, opts.AuthConfig, opts.Logger),
	)

	// Images routes
	mux.HandleFunc("GET /api/v1/images", imagesHandler.ListImages)
	mux.HandleFunc("GET /api/v1/images/{id}", imagesHandler.GetImageByID)
	mux.HandleFunc("GET /api/v1/images/slug/{slug}", imagesHandler.GetImagesBySlug)
	mux.HandleFunc("POST /api/v1/images",
		middleware.NewAuthMiddlewareFunc(imagesHandler.CreateImage, opts.AuthConfig, opts.Logger),
	)
	mux.HandleFunc("PUT /api/v1/images/{id}",
		middleware.NewAuthMiddlewareFunc(imagesHandler.UpdateImage, opts.AuthConfig, opts.Logger),
	)
	mux.HandleFunc("DELETE /api/v1/images/{id}",
		middleware.NewAuthMiddlewareFunc(imagesHandler.DeleteImage, opts.AuthConfig, opts.Logger),
	)

	// Timeline routes
	mux.HandleFunc("GET /api/v1/timelineentries", timelineHandler.ListTimelineEntries)
	mux.HandleFunc("GET /api/v1/timelineentries/{id}", timelineHandler.GetTimelineEntryByID)
	mux.HandleFunc("POST /api/v1/timelineentries",
		middleware.NewAuthMiddlewareFunc(timelineHandler.CreateTimelineEntry, opts.AuthConfig, opts.Logger),
	)
	mux.HandleFunc("PUT /api/v1/timelineentries/{id}",
		middleware.NewAuthMiddlewareFunc(timelineHandler.UpdateTimelineEntry, opts.AuthConfig, opts.Logger),
	)
	mux.HandleFunc("DELETE /api/v1/timelineentries/{id}",
		middleware.NewAuthMiddlewareFunc(timelineHandler.DeleteTimelineEntry, opts.AuthConfig, opts.Logger),
	)

	// Events routes
	mux.HandleFunc("GET /api/v1/events", eventsHandler.GetEvents)

	// GaleryEvent routes
	mux.HandleFunc("GET /api/v1/galery_events", galeryEventHandler.ListGaleryEvents)
	mux.HandleFunc("GET /api/v1/galery_events/{id}", galeryEventHandler.GetGaleryEventByID)
	mux.HandleFunc("POST /api/v1/galery_events",
		middleware.NewAuthMiddlewareFunc(galeryEventHandler.CreateGaleryEvent, opts.AuthConfig, opts.Logger),
	)
	mux.HandleFunc("PUT /api/v1/galery_events",
		middleware.NewAuthMiddlewareFunc(galeryEventHandler.ModifyGaleryEvent, opts.AuthConfig, opts.Logger),
	)
	mux.HandleFunc("DELETE /api/v1/galery_events/{id}",
		middleware.NewAuthMiddlewareFunc(galeryEventHandler.DeleteGaleryEvent, opts.AuthConfig, opts.Logger),
	)

	// Authorization check endpoint (always requires authentication)
	mux.HandleFunc("GET /authorized",
		middleware.NewForceAuthMiddlewareFunc(authHandler.Authorized, opts.AuthConfig, opts.Logger),
	)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	})

	// Apply middleware (outermost to innermost)
	var handler http.Handler = mux
	handler = middleware.Recovery(handler)
	handler = middleware.CORS(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RequestID(handler)

	return handler
}
