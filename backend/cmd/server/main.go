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
	"backend/internal/gateway/gcs"
	httpHandler "backend/internal/http"
	authPlatform "backend/internal/platform/auth"
	firestoreRepo "backend/internal/repository/firestore"
	"backend/internal/server"

	firebaseApp "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
)

func main() {
	ctx := context.Background()

	// Configuration
	port := getEnv("PORT", "8080")

	log.Println("Starting Media CMS Backend...")
	log.Printf("Environment: %s", getEnv("RUNTIME_ENV", "development"))

	// Initialize dependencies
	config := initializeConfig()
	eventsClient := initializeEventsClient()
	gcsGateway := initializeGCSGateway(ctx, config)
	defer gcsGateway.Close()
	objectStore := initializeObjectStore(gcsGateway)
	db := initializeDatabase(ctx, config)
	defer db.Close()
	fbApp := initializeFirebaseApp(ctx, config)
	authClient := initializeAuthClient(ctx, fbApp)
	srv := initializeServer(db, objectStore, eventsClient)
	handler := initializeRouter(ctx, srv, authClient, config)

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
		log.Println("  GET  /api/v1/texts")
		log.Println("  GET  /api/v1/images/{id}")
		log.Println("  GET  /api/v1/timelineentries")

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

// initializeConfig initializes and returns the configuration service
func initializeConfig() configs.ConfigClient {
	config, err := configs.NewConfigService()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}
	return config
}

// initializeEventsClient initializes and returns the events client
func initializeEventsClient() server.GrupyEventsPort {
	return clients.NewEventsClient()
}

// initializeGCSGateway initializes and returns the GCS gateway
func initializeGCSGateway(ctx context.Context, config configs.ConfigClient) *gcs.GCSGateway {
	log.Println("Initializing GCS...")
	gcsGateway, err := gcs.NewGCSGatewayWithProvider(ctx, config)
	if err != nil {
		log.Fatalf("Failed to initialize GCS gateway: %v", err)
	}

	gcsConfig, err := config.GetGCSConfig()
	if err != nil {
		log.Fatalf("Failed to get GCS config: %v", err)
	}

	log.Printf("GCS initialized successfully")
	log.Printf("  Bucket: %s", gcsConfig.BucketName)
	log.Printf("  Project: %s", gcsConfig.ProjectID)
	log.Printf("  Public access: %v", gcsConfig.MakePublic)

	return gcsGateway
}

// initializeObjectStore initializes and returns the object store
func initializeObjectStore(gcsGateway *gcs.GCSGateway) server.ObjectStorePort {
	return clients.NewObjectClient(gcsGateway)
}

// initializeDatabase initializes and returns the Firestore database repository
func initializeDatabase(ctx context.Context, config configs.ConfigClient) *firestoreRepo.DBRepository {
	log.Println("Initializing Firestore...")
	db, err := firestoreRepo.NewDBRepositoryWithProvider(ctx, config)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore database: %v", err)
	}

	fbConfig, err := config.GetFirebaseConfig()
	if err != nil {
		log.Fatalf("Failed to get Firebase config: %v", err)
	}

	collections, err := config.GetCollections()
	if err != nil {
		log.Fatalf("Failed to get collections config: %v", err)
	}

	log.Printf("Firestore initialized successfully")
	log.Printf("  Project: %s", fbConfig.ProjectID)
	log.Printf("  Collections: texts=%s, images=%s, timelines=%s",
		collections.Texts, collections.Images, collections.Timelines)

	return db
}

// initializeFirebaseApp initializes and returns the Firebase app
func initializeFirebaseApp(ctx context.Context, config configs.ConfigClient) *firebaseApp.App {
	log.Println("Initializing Firebase App...")
	fbConfig, err := config.GetFirebaseConfigWithJSONBytes()
	if err != nil {
		log.Fatalf("Failed to get Firebase config: %v", err)
	}

	app, err := clients.NewFirebaseAppClient(ctx, fbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase App: %v", err)
	}

	log.Printf("Firebase App initialized successfully")
	log.Printf("  Project: %s", fbConfig.ProjectID)

	return app
}

// initializeAuthClient initializes and returns the Firebase Auth client
func initializeAuthClient(ctx context.Context, app *firebaseApp.App) *firebaseAuth.Client {
	log.Println("Initializing Firebase Auth client...")
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Auth client: %v", err)
	}

	log.Printf("Firebase Auth client initialized successfully")
	return authClient
}

// initializeServer initializes and returns the server
func initializeServer(db server.DBPort, objectStore server.ObjectStorePort, eventsClient server.GrupyEventsPort) server.Server {
	return server.NewServer(db, objectStore, eventsClient)
}

// initializeRouter initializes and returns the HTTP router
func initializeRouter(ctx context.Context, srv server.Server, authClient *firebaseAuth.Client, config configs.ConfigClient) http.Handler {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	routerOpts := httpHandler.RouterOptions{
		Logger: logger,
	}

	authLevel := config.GetAuthLevel()

	if authClient != nil {
		routerOpts.AuthConfig = authPlatform.AuthConfig{
			Client: authClient,
			Level:  authLevel,
		}
		log.Println("Authentication enabled for protected endpoints")
	} else {
		routerOpts.AuthConfig = authPlatform.AuthConfig{
			Client: nil,
			Level:  authPlatform.AuthOptional,
		}
		log.Println("Authentication disabled or unavailable")
	}

	return httpHandler.NewRouter(ctx, srv, routerOpts)
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
