package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/configs"
	"backend/internal/clients"
	httpHandler "backend/internal/http"
	firestoreRepo "backend/internal/repository/firestore"
	"backend/internal/server"
)

func main() {
	ctx := context.Background()

	// Configuration
	port := getEnv("PORT", "8080")

	log.Println("Starting Media CMS Backend...")
	log.Printf("Environment: %s", getEnv("RUNTIME_ENV", "development"))

	// Load configuration
	config, err := configs.NewConfigService()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize dependencies
	eventsClient := clients.NewEventsClient()

	// Initialize database using the config provider pattern
	var db server.DBPort
	var objectStore server.ObjectStorePort = nil // TODO: Wire up object store later

	// Create DB repository using the config provider
	log.Println("Initializing Firestore...")
	dbRepo, err := firestoreRepo.NewDBRepositoryWithProvider(ctx, config)
	if err != nil {
		log.Printf("Warning: Failed to initialize Firestore: %v", err)
		log.Println("Continuing without database (only /events endpoint will work)")
		db = nil
	} else {
		db = dbRepo
		defer dbRepo.Close()

		// Log successful initialization
		fbConfig, _ := config.GetFirebaseConfig()
		collections, _ := config.GetCollections()
		log.Printf("Firestore initialized successfully")
		log.Printf("  Project: %s", fbConfig.ProjectID)
		log.Printf("  Collections: texts=%s, images=%s, timelines=%s",
			collections.Texts, collections.Images, collections.Timelines)
	}

	// Create unified server
	srv := server.NewServer(db, objectStore, eventsClient)

	// Create HTTP router
	handler := httpHandler.NewRouter(srv)

	// Configure HTTP server
	httpSrv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on port %s", port)
		log.Printf("API available at: http://localhost:%s/api/v1", port)
		log.Println("Available endpoints:")
		log.Println("  GET  /api/v1/events")
		if db != nil {
			log.Println("  GET  /api/v1/texts")
			log.Println("  GET  /api/v1/images/{id}")
			log.Println("  GET  /api/v1/timelineentries")
		}

		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 30 second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
