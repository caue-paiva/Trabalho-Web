# Google Cloud Storage Gateway

This package provides the GCS (Google Cloud Storage) adapter implementation for the Media CMS object storage layer.

## Quick Start

```go
import "backend/internal/gateway/gcs"

// Initialize GCS gateway with config
gateway, err := gcs.NewGCSGateway(ctx, gcsConfig)
if err != nil {
    log.Fatal(err)
}

// Use the gateway (implements server.ObjectStorePort)
url, err := gateway.PutObject(ctx, "images/photo.jpg", imageData)
```

## Files

- **`CLAUDE.md`** - Complete implementation specification (READ THIS FIRST)
- **`storage.go`** - Main GCS gateway implementation
- **`storage_test.go`** - Integration tests with real GCS bucket
- **`config.go`** - Configuration types and provider pattern
- **`helpers.go`** - Content-type detection and utilities

## Interface

Implements `server.ObjectStorePort`:
- `PutObject(ctx, key, data) -> (url, error)` - Upload file
- `DeleteObject(ctx, key) -> error` - Delete file
- `SignedURL(ctx, key) -> (url, error)` - Generate temporary signed URL

## Configuration

Add to `configs/development.yaml`:

```yaml
gcs:
  bucket_name: sitegrupysanca.appspot.com
  credentials_path: sitegrupysanca-firebase-adminsdk-fbsvc-ff7567bd6e.json
  project_id: sitegrupysanca
  make_public: true
  signed_url_expiry_minutes: 15
```

## Dependencies

```bash
go get cloud.google.com/go/storage
go get google.golang.org/api/option
```

## Testing

Run integration tests (requires GCS access):

```bash
go test -v ./internal/gateway/gcs
```

## Documentation

See **`CLAUDE.md`** for:
- Detailed architecture
- Configuration integration
- Implementation requirements
- Security considerations
- Testing strategy
- Migration guide from mock

## Status

🚧 **Not yet implemented** - Currently using mock object store (`clients.NewMockObjectStore()`)

Implementation checklist in `CLAUDE.md` section 15.
