package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/clients"
	httpHandler "backend/internal/http"
	"backend/internal/service"
)

func main() {
	// Configuration
	port := getEnv("PORT", "8080")

	log.Println("Starting Media CMS Backend...")

	// Initialize dependencies
	eventsClient := clients.NewEventsClient()

	// Note: DB and ObjectStore are nil since we haven't wired them up yet
	// Only the /events endpoint will work for now
	var db service.DBPort = nil
	var objectStore service.ObjectStorePort = nil

	// Create unified server
	server := service.NewServer(db, objectStore, eventsClient)

	// Create HTTP router
	handler := httpHandler.NewRouter(server)

	// Configure server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on port %s", port)
		log.Printf("Events API available at: http://localhost:%s/api/v1/events", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
