# Google Cloud Storage Gateway Implementation

This document specifies the implementation for the Google Cloud Storage (GCS) gateway adapter that satisfies the `ObjectStorePort` interface defined in `internal/server/ports.go`.

---

## 1. Overview

The GCS gateway provides concrete implementation of object storage operations using Google Cloud Storage, which is the underlying storage for Firebase Storage. This adapter allows the application to:

- Upload images and media files to GCS
- Delete objects from GCS
- Generate signed URLs for secure, temporary access to private objects
- Integrate seamlessly with Firebase Storage buckets

---

## 2. Architecture Position

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Handler Layer                    │
│                 (receives image upload)                  │
└────────────────────────┬────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                    Service Layer                         │
│              (internal/server/image.go)                  │
│            UploadImage(ctx, meta, data)                  │
└────────────────────────┬────────────────────────────────┘
                         │
                         ▼ calls ObjectStorePort
┌─────────────────────────────────────────────────────────┐
│                  ObjectStorePort                         │
│            PutObject(ctx, key, data)                     │
│            DeleteObject(ctx, key)                        │
│            SignedURL(ctx, key)                           │
└────────────────────────┬────────────────────────────────┘
                         │
                         ▼ implements
┌─────────────────────────────────────────────────────────┐
│              GCS Gateway (this module)                   │
│          internal/gateway/gcs/storage.go                 │
│                                                          │
│   Uses: cloud.google.com/go/storage                     │
│   Bucket: sitegrupysanca.appspot.com                    │
└──────────────────────────────────────────────────────────┘
```

---

## 3. Interface Contract

The implementation must satisfy this interface from `internal/server/ports.go`:

```go
type ObjectStorePort interface {
    PutObject(ctx context.Context, key string, data []byte) (publicURL string, err error)
    DeleteObject(ctx context.Context, key string) error
    SignedURL(ctx context.Context, key string) (string, error)
}
```

### Method Specifications

#### `PutObject(ctx, key, data) -> (publicURL, error)`
- **Purpose**: Upload a file to GCS
- **Input**:
  - `key`: Object path (e.g., "images/2024/photo-abc123.jpg")
  - `data`: Raw file bytes
- **Output**: Public URL to access the object
- **Behavior**:
  - Set appropriate `Content-Type` based on file extension
  - Make object publicly readable OR return signed URL
  - Handle upload failures gracefully
  - Return URL in format: `https://storage.googleapis.com/bucket-name/key`

#### `DeleteObject(ctx, key) -> error`
- **Purpose**: Delete an object from GCS
- **Input**: `key` - Object path to delete
- **Behavior**:
  - If object doesn't exist, return nil (idempotent)
  - Handle permission errors appropriately
  - Log deletion for audit trail

#### `SignedURL(ctx, key) -> (string, error)`
- **Purpose**: Generate temporary signed URL for private objects
- **Input**: `key` - Object path
- **Output**: Time-limited signed URL (e.g., valid for 15 minutes)
- **Use case**: Securely share private images without making them public

---

## 4. Configuration Integration

### 4.1 Configuration Structure

Add to `configs/development.yaml` and `configs/production.yaml`:

```yaml
# Google Cloud Storage configuration
gcs:
  bucket_name: sitegrupysanca.appspot.com
  credentials_path: sitegrupysanca-firebase-adminsdk-fbsvc-ff7567bd6e.json
  project_id: sitegrupysanca
  make_public: true  # Whether uploaded objects should be publicly accessible
  signed_url_expiry_minutes: 15  # For private objects, how long signed URLs last
```

### 4.2 Config Client Extension

Extend `configs/config.go` to support GCS configuration:

```go
// GCSConfig holds Google Cloud Storage configuration
type GCSConfig struct {
    BucketName              string `yaml:"bucket_name"`
    CredentialsPath         string `yaml:"credentials_path"`
    ProjectID               string `yaml:"project_id"`
    MakePublic              bool   `yaml:"make_public"`
    SignedURLExpiryMinutes  int    `yaml:"signed_url_expiry_minutes"`
}
```

Add method to `ConfigClient` interface:

```go
type ConfigClient interface {
    // ... existing methods ...

    // GetGCSConfig returns the Google Cloud Storage configuration
    GetGCSConfig() (GCSConfig, error)
}
```

Implement in `configService`:

```go
func (s *configService) GetGCSConfig() (GCSConfig, error) {
    var config GCSConfig
    if err := s.UnmarshalKey("gcs", &config); err != nil {
        return GCSConfig{}, err
    }
    return config, nil
}
```

---

## 5. Client Initialization

### 5.1 Gateway Constructor

The GCS gateway should provide a constructor that accepts configuration:

```go
// Expected signature (implementation not shown)
func NewGCSGateway(ctx context.Context, config GCSConfig) (server.ObjectStorePort, error)
```

**Responsibilities**:
1. Read credentials from file path specified in config
2. Create `storage.Client` using credentials
3. Validate bucket exists and is accessible
4. Store bucket reference for operations
5. Return implementation of `ObjectStorePort`

### 5.2 Authentication Methods

Support multiple authentication methods (in order of precedence):

1. **Service Account JSON** (production)
   - Use `option.WithCredentialsFile(path)` or `option.WithCredentialsJSON(bytes)`
   - Read from `config.CredentialsPath`

2. **Application Default Credentials** (development)
   - Use `storage.NewClient(ctx)` with no options
   - Works when running on GCP (Cloud Run, GCE) or with `gcloud auth application-default login`

3. **Fallback to Mock** (testing)
   - If no credentials available and not in production, use mock implementation

---

## 6. Integration with Main Application

### 6.1 Dependency Wiring in `cmd/server/main.go`

Replace the mock object store with real GCS implementation:

```go
// Before (current):
objectStore := clients.NewMockObjectStore()
log.Println("Using mock object store (no actual file storage)")

// After (with GCS):
gcsConfig, err := config.GetGCSConfig()
if err != nil {
    log.Printf("Warning: Failed to load GCS config: %v", err)
    objectStore = clients.NewMockObjectStore() // Fallback to mock
} else {
    objectStore, err = gcsGateway.NewGCSGateway(ctx, gcsConfig)
    if err != nil {
        log.Fatalf("Failed to initialize GCS gateway: %v", err)
    }
    log.Printf("GCS initialized successfully (bucket: %s)", gcsConfig.BucketName)
}
```

### 6.2 Provider Pattern (Optional Enhancement)

For consistency with Firestore, could create a provider interface:

```go
// In internal/gateway/gcs/config.go
type GCSConfigProvider interface {
    GetGCSConfig() (configs.GCSConfig, error)
    GetCredentialsJSON(filename string) ([]byte, error)
}

func NewGCSGatewayWithProvider(ctx context.Context, provider GCSConfigProvider) (server.ObjectStorePort, error)
```

Then in main.go:

```go
objectStore, err := gcsGateway.NewGCSGatewayWithProvider(ctx, config)
```

---

## 7. File Structure

```
/internal/gateway/gcs/
├── CLAUDE.md              # This file - implementation specification
├── storage.go             # Main GCS gateway implementation
├── storage_test.go        # Unit tests with real GCS integration
├── config.go              # Config types and provider pattern (optional)
└── helpers.go             # Helper functions (content-type detection, URL generation)
```

---

## 8. Key Implementation Details

### 8.1 Content-Type Detection

The gateway should set appropriate `Content-Type` headers:

```
.jpg, .jpeg -> image/jpeg
.png        -> image/png
.gif        -> image/gif
.webp       -> image/webp
.svg        -> image/svg+xml
```

Use Go's `mime.TypeByExtension()` or explicit mapping.

### 8.2 Object Naming Convention

The service layer generates keys like:
- `images/{slug}-{timestamp}` for user uploads
- `processed/{slug}-{size}-{timestamp}` for thumbnails (future)

The gateway should:
- Accept the key as-is (no modification)
- Create parent "folders" automatically (GCS handles this)

### 8.3 Public vs Private Objects

**Option A: Public by default** (simpler)
- Set ACL to public-read on upload
- Return standard GCS public URL: `https://storage.googleapis.com/bucket/key`

**Option B: Private with signed URLs** (more secure)
- Keep objects private
- Generate signed URLs with expiry for access
- Requires managing credentials for signing

Recommendation: Start with **Option A** (public), add Option B later if needed.

### 8.4 Error Handling

Map GCS errors to application errors:
- `storage.ErrObjectNotExist` → "object not found"
- Permission errors → "access denied"
- Network errors → "storage unavailable"

---

## 9. Testing Strategy

### 9.1 Real Integration Tests

Similar to Firestore tests, create `storage_test.go` that:
1. Uses real GCS bucket (test bucket or dev bucket)
2. Uploads test files
3. Verifies URLs are accessible
4. Deletes test files (cleanup)
5. Tests signed URL generation

### 9.2 Mock for Unit Tests

Keep `clients/mock_objectstore.go` for:
- Service layer unit tests
- HTTP handler tests
- Local development without GCS access

---

## 10. Security Considerations

### 10.1 Credentials Management

- **Never commit** service account JSON to git
- Store credentials file path in config, not the actual JSON
- Use environment-specific credentials (dev vs prod)
- Consider using Secret Manager for production

### 10.2 Bucket Security

- Configure Firebase Storage security rules
- Limit public access to specific paths if needed
- Use Cloud IAM for service account permissions

### 10.3 Content Validation

While the gateway doesn't validate content, the service layer should:
- Validate file size (current: 10MB limit in `image.go`)
- Validate content type (only allow images)
- Sanitize filenames/keys

---

## 11. Future Enhancements

### 11.1 Image Manipulation Pipeline

Using Cloud Functions or in-process libraries:
1. Upload original to `images/original/`
2. Trigger processing (resize, optimize)
3. Store variants in `images/thumb/`, `images/medium/`
4. Update DB with all URLs

### 11.2 CDN Integration

- Enable Cloud CDN for the bucket
- Add `Cache-Control` headers on upload
- Return CDN URLs instead of storage URLs

### 11.3 Multipart Uploads

For large files (videos), implement resumable uploads:
- Use `wc.ChunkSize` to control chunk size
- Handle retry logic for failed chunks

---

## 12. Dependencies

Add to `go.mod`:

```bash
go get cloud.google.com/go/storage
go get google.golang.org/api/option
```

These are the official Google Cloud Go client libraries and are well-maintained.

---

## 13. Example Usage Flow

1. **User uploads image via API** → `POST /api/v1/images` with base64 data
2. **HTTP handler** decodes base64 → calls service layer
3. **Service layer** (`image.go`):
   - Generates object key: `images/gallery-1-1699123456`
   - Validates file size (10MB limit)
   - Calls `objectStore.PutObject(ctx, key, data)`
4. **GCS Gateway**:
   - Uploads to `gs://sitegrupysanca.appspot.com/images/gallery-1-1699123456`
   - Returns URL: `https://storage.googleapis.com/sitegrupysanca.appspot.com/images/gallery-1-1699123456`
5. **Service layer**:
   - Updates `meta.ObjectURL` with returned URL
   - Saves metadata to Firestore
6. **Client receives** JSON with `object_url` field pointing to GCS

---

## 14. Migration from Mock

Current state:
```go
objectStore := clients.NewMockObjectStore()
```

After implementation:
```go
objectStore, err := gcsGateway.NewGCSGateway(ctx, gcsConfig)
if err != nil {
    log.Fatalf("Failed to initialize GCS: %v", err)
}
defer objectStore.Close() // If GCS gateway needs cleanup
```

All existing image upload/delete flows will work automatically since they use the `ObjectStorePort` interface!

---

## 15. Checklist for Implementation

- [ ] Add GCS configuration to `development.yaml` and `production.yaml`
- [ ] Extend `configs/config.go` with `GCSConfig` struct and `GetGCSConfig()` method
- [ ] Create `internal/gateway/gcs/storage.go` with gateway implementation
- [ ] Implement `PutObject()` with public URL return
- [ ] Implement `DeleteObject()` with idempotent behavior
- [ ] Implement `SignedURL()` for future use
- [ ] Add content-type detection helper
- [ ] Create `storage_test.go` with integration tests
- [ ] Update `cmd/server/main.go` to wire up GCS gateway
- [ ] Test full upload flow via API
- [ ] Test delete flow
- [ ] Verify URLs are accessible from browser
- [ ] Update documentation in main CLAUDE.md

---

## References

- [Google Cloud Storage Go Client](https://pkg.go.dev/cloud.google.com/go/storage)
- [Firebase Storage Documentation](https://firebase.google.com/docs/storage)
- [GCS Access Control](https://cloud.google.com/storage/docs/access-control)
- [Signed URLs Documentation](https://cloud.google.com/storage/docs/access-control/signed-urls)
