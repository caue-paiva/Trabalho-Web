package clients

import (
	"context"

	"backend/internal/server"
)

// ObjectStoreGateway defines the interface for the underlying storage gateway
// This allows the client to wrap any gateway implementation (GCS, S3, etc.)
type ObjectStoreGateway interface {
	PutObject(ctx context.Context, key string, data []byte) (string, error)
	DeleteObject(ctx context.Context, key string) error
	SignedURL(ctx context.Context, key string) (string, error)
	Close() error
}

// Compile-time interface check
var _ server.ObjectStorePort = (*objectClient)(nil)

type objectClient struct {
	gateway ObjectStoreGateway
}

// NewObjectClient creates a new ObjectStorePort implementation that wraps a gateway
func NewObjectClient(gateway ObjectStoreGateway) server.ObjectStorePort {
	return &objectClient{
		gateway: gateway,
	}
}

// PutObject uploads an object via the gateway
func (c *objectClient) PutObject(ctx context.Context, key string, data []byte) (string, error) {
	return c.gateway.PutObject(ctx, key, data)
}

// DeleteObject deletes an object via the gateway
func (c *objectClient) DeleteObject(ctx context.Context, key string) error {
	return c.gateway.DeleteObject(ctx, key)
}

// SignedURL generates a signed URL via the gateway
func (c *objectClient) SignedURL(ctx context.Context, key string) (string, error) {
	return c.gateway.SignedURL(ctx, key)
}

// Close closes the underlying gateway connection
func (c *objectClient) Close() error {
	if c.gateway != nil {
		return c.gateway.Close()
	}
	return nil
}
