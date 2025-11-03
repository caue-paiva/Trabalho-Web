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

	// Load Firebase configuration from config file
	var db server.DBPort
	var objectStore server.ObjectStorePort = nil // TODO: Wire up object store later

	// Unmarshal Firebase config
	type FirebaseConfig struct {
		ProjectID       string `yaml:"project_id"`
		CredentialsPath string `yaml:"credentials_path"`
	}

	var fbConfig FirebaseConfig
	if err := config.UnmarshalKey("firebase", &fbConfig); err != nil {
		log.Printf("Warning: Failed to load Firebase config: %v", err)
		log.Println("Continuing without database (only /events endpoint will work)")
		db = nil
	} else {
		// Unmarshal collection names
		type Collections struct {
			Texts     string `yaml:"texts"`
			Images    string `yaml:"images"`
			Timelines string `yaml:"timelines"`
		}

		var cols Collections
		if err := config.UnmarshalKey("collections", &cols); err != nil {
			log.Fatalf("Failed to load collection names: %v", err)
		}

		// Create Firestore collection names
		collections := firestoreRepo.CollectionNames{
			Texts:           cols.Texts,
			Images:          cols.Images,
			TimelineEntries: cols.Timelines,
		}

		// Configure Firestore client
		firestoreConfig := firestoreRepo.FirestoreConfig{
			ProjectID:       fbConfig.ProjectID,
			CredentialsPath: fbConfig.CredentialsPath,
			Collections:     collections,
		}

		// Initialize Firestore client
		log.Printf("Initializing Firestore (project: %s)...", fbConfig.ProjectID)
		firestoreClient, err := firestoreRepo.NewFirestoreClient(ctx, firestoreConfig)
		if err != nil {
			log.Fatalf("Failed to create Firestore client: %v", err)
		}
		defer firestoreClient.Close()

		// Create DB repository
		db = firestoreRepo.NewDBRepository(firestoreClient, collections)
		log.Printf("Firestore initialized successfully")
		log.Printf("  Collections: texts=%s, images=%s, timelines=%s",
			collections.Texts, collections.Images, collections.TimelineEntries)
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
