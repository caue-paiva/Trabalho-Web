package clients

import (
	"context"
	"fmt"

	"backend/internal/server"
)

// Compile-time check that mockObjectStore implements server.ObjectStorePort
var _ server.ObjectStorePort = (*mockObjectStore)(nil)

// mockObjectStore is a no-op implementation of ObjectStorePort for testing
type mockObjectStore struct{}

// NewMockObjectStore creates a new mock object store that doesn't actually store anything
// It returns fake URLs and does no-op operations, useful for testing without real storage
func NewMockObjectStore() server.ObjectStorePort {
	return &mockObjectStore{}
}

// PutObject returns a mock URL without actually storing data
func (m *mockObjectStore) PutObject(ctx context.Context, key string, data []byte) (publicURL string, err error) {
	// Return a fake URL that includes the key for debugging
	mockURL := fmt.Sprintf("https://mock-storage.example.com/%s", key)
	return mockURL, nil
}

// DeleteObject is a no-op, always succeeds
func (m *mockObjectStore) DeleteObject(ctx context.Context, key string) error {
	// No-op: pretend we deleted it
	return nil
}

// SignedURL returns a fake signed URL
func (m *mockObjectStore) SignedURL(ctx context.Context, key string) (string, error) {
	// Return a fake signed URL
	mockSignedURL := fmt.Sprintf("https://mock-storage.example.com/%s?signed=true", key)
	return mockSignedURL, nil
}
