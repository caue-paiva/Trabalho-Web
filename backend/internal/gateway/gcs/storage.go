package gcs

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"backend/configs"
)

const (
	_defaultExpiryInMinutes = 15
	_publicReadACL          = "publicRead"
	_cacheControlImmutable  = "public, max-age=31536000" // 1 year cache for immutable content
)

// GCSGateway implements object storage operations using Google Cloud Storage
type GCSGateway struct {
	client                 *storage.Client
	bucket                 *storage.BucketHandle
	bucketName             string
	makePublic             bool
	signedURLExpiryMinutes int
}

// NewGCSGateway creates a new GCS gateway with the given configuration
func NewGCSGateway(ctx context.Context, config configs.GCSConfig) (*GCSGateway, error) {
	// Validate configuration
	if config.BucketName == "" {
		return nil, fmt.Errorf("bucket_name is required in GCS configuration")
	}
	if config.ProjectID == "" {
		return nil, fmt.Errorf("project_id is required in GCS configuration")
	}

	// Create storage client with credentials
	var client *storage.Client
	var err error

	if config.CredentialsPath != "" {
		// Use service account credentials from file
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(config.CredentialsPath))
		if err != nil {
			return nil, fmt.Errorf("failed to create GCS client with credentials file: %w", err)
		}
	} else {
		// Use application default credentials (for Cloud Run, GCE, or gcloud auth)
		client, err = storage.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create GCS client with default credentials: %w", err)
		}
	}

	// Get bucket handle
	bucket := client.Bucket(config.BucketName)

	// Verify bucket exists and is accessible
	_, err = bucket.Attrs(ctx)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to access bucket %s: %w (verify bucket exists and credentials have access)", config.BucketName, err)
	}

	// Set default expiry if not specified
	expiryMinutes := config.SignedURLExpiryMinutes
	if expiryMinutes == 0 {
		expiryMinutes = _defaultExpiryInMinutes
	}

	return &GCSGateway{
		client:                 client,
		bucket:                 bucket,
		bucketName:             config.BucketName,
		makePublic:             config.MakePublic,
		signedURLExpiryMinutes: expiryMinutes,
	}, nil
}

// NewGCSGatewayWithProvider creates a new GCS gateway using a config provider
// This follows the same pattern as Firestore for consistency
func NewGCSGatewayWithProvider(ctx context.Context, provider configs.ConfigClient) (*GCSGateway, error) {
	config, err := provider.GetGCSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get GCS config: %w", err)
	}
	return NewGCSGateway(ctx, config)
}

// PutObject uploads a file to GCS and returns its public URL
func (g *GCSGateway) PutObject(ctx context.Context, key string, data []byte) (string, error) {
	// Create object writer
	obj := g.bucket.Object(key)
	writer := obj.NewWriter(ctx)

	// Set content type based on file extension
	writer.ContentType = detectContentType(key)

	// Set cache control for long-term caching (immutable content)
	writer.CacheControl = _cacheControlImmutable

	// Set ACL to public during upload if configured (no separate network call needed)
	if g.makePublic {
		writer.PredefinedACL = _publicReadACL
	}

	// Write data
	if _, err := writer.Write(data); err != nil {
		writer.Close()
		return "", fmt.Errorf("failed to write object data: %w", err)
	}

	// Close writer (this commits the upload)
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close object writer: %w", err)
	}

	// Return public URL
	publicURL := g.getPublicURL(key)
	return publicURL, nil
}

// DeleteObject deletes an object from GCS
// Returns nil if object doesn't exist (idempotent operation)
func (g *GCSGateway) DeleteObject(ctx context.Context, key string) error {
	obj := g.bucket.Object(key)

	// Delete the object
	if err := obj.Delete(ctx); err != nil {
		// Check if error is because object doesn't exist
		if err == storage.ErrObjectNotExist {
			// Idempotent - return success if object doesn't exist
			return nil
		}
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// SignedURL generates a temporary signed URL for private object access
func (g *GCSGateway) SignedURL(ctx context.Context, key string) (string, error) {
	// Calculate expiry time
	expires := time.Now().Add(time.Duration(g.signedURLExpiryMinutes) * time.Minute)

	// Generate signed URL
	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: expires,
	}

	url, err := g.bucket.SignedURL(key, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// getPublicURL constructs the public URL for an object
func (g *GCSGateway) getPublicURL(key string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, key)
}

// Close closes the GCS client connection
func (g *GCSGateway) Close() error {
	if g.client != nil {
		return g.client.Close()
	}
	return nil
}

// GetObject retrieves an object's content from GCS (helper method, not in port interface)
// This can be useful for testing or future features
func (g *GCSGateway) GetObject(ctx context.Context, key string) ([]byte, error) {
	obj := g.bucket.Object(key)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, fmt.Errorf("object not found: %s", key)
		}
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	return data, nil
}
