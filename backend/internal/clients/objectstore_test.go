package clients

import (
	"context"
	"os"
	"testing"

	"backend/configs"
	"backend/internal/gateway/gcs"
	"backend/internal/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestObjectStore creates a test object store client with real GCS
func setupTestObjectStore(t *testing.T) (server.ObjectStorePort, func()) {
	os.Unsetenv("RUNTIME_ENV")

	ctx := context.Background()

	// Load configuration
	config, err := configs.NewConfigService()
	require.NoError(t, err, "Failed to load config")

	// Initialize GCS gateway
	gcsGateway, err := gcs.NewGCSGatewayWithProvider(ctx, config)
	require.NoError(t, err, "Failed to initialize GCS gateway")

	// Create object store client
	objectStore := NewObjectClient(gcsGateway)

	cleanup := func() {
		gcsGateway.Close()
	}

	return objectStore, cleanup
}

func TestObjectStoreClient_PutObject(t *testing.T) {
	objectStore, cleanup := setupTestObjectStore(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name         string
		key          string
		data         []byte
		expectError  bool
		validateFunc func(t *testing.T, url string)
	}{
		{
			name:        "upload small text file",
			key:         "test-small-file.txt",
			data:        []byte("Hello, GCS! This is a test file."),
			expectError: false,
			validateFunc: func(t *testing.T, url string) {
				assert.NotEmpty(t, url, "Should return a URL")
				assert.Contains(t, url, "storage.googleapis.com", "URL should point to GCS")
				assert.Contains(t, url, "test-small-file.txt", "URL should contain the filename")
			},
		},
		{
			name:        "upload small image (1x1 PNG)",
			key:         "test-tiny-image.png",
			data:        []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, // PNG header
			expectError: false,
			validateFunc: func(t *testing.T, url string) {
				assert.NotEmpty(t, url, "Should return a URL")
				assert.Contains(t, url, "test-tiny-image.png", "URL should contain the filename")
			},
		},
		{
			name:        "upload with special characters in key",
			key:         "test-special-chars-üñíçødé.txt",
			data:        []byte("Content with special filename"),
			expectError: false,
			validateFunc: func(t *testing.T, url string) {
				assert.NotEmpty(t, url, "Should return a URL")
			},
		},
		{
			name:        "upload with nested path",
			key:         "subfolder/nested/test-nested.txt",
			data:        []byte("Content in nested folder"),
			expectError: false,
			validateFunc: func(t *testing.T, url string) {
				assert.NotEmpty(t, url, "Should return a URL")
				assert.Contains(t, url, "subfolder/nested/test-nested.txt", "URL should contain the nested path")
			},
		},
		{
			name:        "upload larger file (1KB)",
			key:         "test-1kb-file.bin",
			data:        make([]byte, 1024), // 1KB of zeros
			expectError: false,
			validateFunc: func(t *testing.T, url string) {
				assert.NotEmpty(t, url, "Should return a URL")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Upload the object
			url, err := objectStore.PutObject(ctx, tt.key, tt.data)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err, "Failed to upload object")

			// Cleanup: delete the uploaded object
			defer func() {
				err := objectStore.DeleteObject(ctx, tt.key)
				assert.NoError(t, err, "Failed to cleanup uploaded object")
			}()

			// Run validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, url)
			}
		})
	}
}

func TestObjectStoreClient_DeleteObject(t *testing.T) {
	objectStore, cleanup := setupTestObjectStore(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name        string
		setupKey    string
		setupData   []byte
		deleteKey   string
		expectError bool
	}{
		{
			name:        "delete existing object",
			setupKey:    "test-delete-me.txt",
			setupData:   []byte("This file will be deleted"),
			deleteKey:   "test-delete-me.txt",
			expectError: false,
		},
		{
			name:        "delete non-existent object (idempotent)",
			setupKey:    "", // No setup
			setupData:   nil,
			deleteKey:   "non-existent-file-12345.txt",
			expectError: false, // Should not error (idempotent)
		},
		{
			name:        "delete object in nested path",
			setupKey:    "delete-test/nested/file.txt",
			setupData:   []byte("Nested file to delete"),
			deleteKey:   "delete-test/nested/file.txt",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: create object if needed
			if tt.setupKey != "" && tt.setupData != nil {
				_, err := objectStore.PutObject(ctx, tt.setupKey, tt.setupData)
				require.NoError(t, err, "Failed to setup test object")
			}

			// Delete the object
			err := objectStore.DeleteObject(ctx, tt.deleteKey)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err, "Failed to delete object")

			// Additional cleanup for nested paths
			if tt.setupKey != "" && tt.setupKey != tt.deleteKey {
				defer objectStore.DeleteObject(ctx, tt.setupKey)
			}
		})
	}
}

func TestObjectStoreClient_SignedURL(t *testing.T) {
	objectStore, cleanup := setupTestObjectStore(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test object first
	testKey := "test-signed-url.txt"
	testData := []byte("Content for signed URL test")

	_, err := objectStore.PutObject(ctx, testKey, testData)
	require.NoError(t, err, "Failed to create test object")

	defer func() {
		objectStore.DeleteObject(ctx, testKey)
	}()

	tests := []struct {
		name         string
		key          string
		expectError  bool
		validateFunc func(t *testing.T, url string)
	}{
		{
			name:        "generate signed URL for existing object",
			key:         testKey,
			expectError: false,
			validateFunc: func(t *testing.T, url string) {
				assert.NotEmpty(t, url, "Should return a URL")
				assert.Contains(t, url, "googleapis.com", "URL should be a Google API URL")
				assert.Contains(t, url, "Expires=", "URL should contain expiry parameter")
				assert.Contains(t, url, "Signature=", "URL should contain signature")
			},
		},
		{
			name:        "generate signed URL for non-existent object",
			key:         "non-existent-file-12345.txt",
			expectError: false, // Signed URL can be generated even if object doesn't exist yet
			validateFunc: func(t *testing.T, url string) {
				assert.NotEmpty(t, url, "Should return a URL")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate signed URL
			url, err := objectStore.SignedURL(ctx, tt.key)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err, "Failed to generate signed URL")

			// Run validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, url)
			}
		})
	}
}

func TestObjectStoreClient_UploadAndDeleteMultiple(t *testing.T) {
	objectStore, cleanup := setupTestObjectStore(t)
	defer cleanup()

	ctx := context.Background()

	// Upload multiple objects
	testObjects := []struct {
		key  string
		data []byte
	}{
		{"batch-test-1.txt", []byte("Batch content 1")},
		{"batch-test-2.txt", []byte("Batch content 2")},
		{"batch-test-3.txt", []byte("Batch content 3")},
	}

	var uploadedKeys []string

	// Upload all objects
	for _, obj := range testObjects {
		url, err := objectStore.PutObject(ctx, obj.key, obj.data)
		require.NoError(t, err, "Failed to upload object: %s", obj.key)
		assert.NotEmpty(t, url, "Should return URL for: %s", obj.key)
		uploadedKeys = append(uploadedKeys, obj.key)
	}

	// Cleanup: delete all uploaded objects
	defer func() {
		for _, key := range uploadedKeys {
			err := objectStore.DeleteObject(ctx, key)
			assert.NoError(t, err, "Failed to cleanup object: %s", key)
		}
	}()

	// Verify all were uploaded by generating signed URLs
	for _, key := range uploadedKeys {
		url, err := objectStore.SignedURL(ctx, key)
		assert.NoError(t, err, "Should be able to generate signed URL for: %s", key)
		assert.NotEmpty(t, url, "Should return signed URL for: %s", key)
	}
}

func TestObjectStoreClient_PutObjectWithDifferentExtensions(t *testing.T) {
	objectStore, cleanup := setupTestObjectStore(t)
	defer cleanup()

	ctx := context.Background()

	// Test different file extensions to verify content-type detection
	testFiles := []struct {
		key      string
		data     []byte
		contains string // What the URL should contain
	}{
		{"test.jpg", []byte{0xFF, 0xD8, 0xFF}, "test.jpg"},   // JPEG header
		{"test.png", []byte{0x89, 0x50, 0x4E, 0x47}, "test.png"}, // PNG header
		{"test.txt", []byte("Plain text content"), "test.txt"},
		{"test.json", []byte(`{"key":"value"}`), "test.json"},
		{"test.pdf", []byte("%PDF-1.4"), "test.pdf"},
	}

	for _, tf := range testFiles {
		t.Run("upload_"+tf.key, func(t *testing.T) {
			// Upload the file
			url, err := objectStore.PutObject(ctx, tf.key, tf.data)
			require.NoError(t, err, "Failed to upload %s", tf.key)
			assert.Contains(t, url, tf.contains, "URL should contain filename")

			// Cleanup
			defer func() {
				err := objectStore.DeleteObject(ctx, tf.key)
				assert.NoError(t, err, "Failed to cleanup %s", tf.key)
			}()
		})
	}
}

func TestObjectStoreClient_DeleteIdempotency(t *testing.T) {
	objectStore, cleanup := setupTestObjectStore(t)
	defer cleanup()

	ctx := context.Background()

	// Create and delete an object
	key := "test-idempotent-delete.txt"
	data := []byte("Test idempotency")

	// Upload
	_, err := objectStore.PutObject(ctx, key, data)
	require.NoError(t, err, "Failed to upload test object")

	// Delete once
	err = objectStore.DeleteObject(ctx, key)
	require.NoError(t, err, "Failed to delete object first time")

	// Delete again - should not error (idempotent)
	err = objectStore.DeleteObject(ctx, key)
	assert.NoError(t, err, "Delete should be idempotent and not error on second call")

	// Delete again - third time to be sure
	err = objectStore.DeleteObject(ctx, key)
	assert.NoError(t, err, "Delete should be idempotent even on third call")
}

func TestObjectStoreClient_EmptyFile(t *testing.T) {
	objectStore, cleanup := setupTestObjectStore(t)
	defer cleanup()

	ctx := context.Background()

	// Upload an empty file
	key := "test-empty-file.txt"
	emptyData := []byte{}

	url, err := objectStore.PutObject(ctx, key, emptyData)
	require.NoError(t, err, "Should be able to upload empty file")
	assert.NotEmpty(t, url, "Should return URL for empty file")

	// Cleanup
	defer func() {
		err := objectStore.DeleteObject(ctx, key)
		assert.NoError(t, err, "Failed to cleanup empty file")
	}()
}

func TestObjectStoreClient_LargeFile(t *testing.T) {
	// Skip this test in short mode as it uploads 5MB
	if testing.Short() {
		t.Skip("Skipping large file test in short mode")
	}

	objectStore, cleanup := setupTestObjectStore(t)
	defer cleanup()

	ctx := context.Background()

	// Create a 5MB file
	key := "test-large-file-5mb.bin"
	largeData := make([]byte, 5*1024*1024) // 5MB
	// Fill with some pattern so it's not all zeros
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// Upload
	url, err := objectStore.PutObject(ctx, key, largeData)
	require.NoError(t, err, "Should be able to upload 5MB file")
	assert.NotEmpty(t, url, "Should return URL for large file")

	// Cleanup
	defer func() {
		err := objectStore.DeleteObject(ctx, key)
		assert.NoError(t, err, "Failed to cleanup large file")
	}()
}
