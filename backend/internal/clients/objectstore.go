package clients

import (
	"context"
	"fmt"

	"backend/internal/server"
)

// Compile-time interface check
var _ server.ObjectStorePort = (*objectClient)(nil)

type objectClient struct {
	// Future: gateway dependencies
}

// NewObjectClient creates a new ObjectStorePort implementation (stub)
func NewObjectClient() server.ObjectStorePort {
	return &objectClient{}
}

func (c *objectClient) PutObject(ctx context.Context, key string, data []byte) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *objectClient) DeleteObject(ctx context.Context, key string) error {
	return fmt.Errorf("not implemented")
}

func (c *objectClient) SignedURL(ctx context.Context, key string) (string, error) {
	return "", fmt.Errorf("not implemented")
}
