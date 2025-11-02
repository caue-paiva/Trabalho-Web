package clients

import (
	"context"
	"fmt"

	"backend/internal/service"
)

// Compile-time interface check
var _ service.ObjectStorePort = (*objectClient)(nil)

type objectClient struct {
	// Future: gateway dependencies
}

// NewObjectClient creates a new ObjectStorePort implementation (stub)
func NewObjectClient() service.ObjectStorePort {
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
