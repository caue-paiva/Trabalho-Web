# Object Store Client Integration Tests

This directory contains integration tests for the object store client that test against real Google Cloud Storage.

## Test File

- **`objectstore_test.go`** - Integration tests for GCS object storage operations

## Running the Tests

### Run all object store tests:
```bash
go test ./internal/clients -v -run TestObjectStoreClient
```

### Run specific test:
```bash
go test ./internal/clients -v -run TestObjectStoreClient_PutObject
```

### Run with short mode (skips large file tests):
```bash
go test ./internal/clients -v -short -run TestObjectStoreClient
```

## Test Coverage

### Upload Tests (`TestObjectStoreClient_PutObject`)
- Upload small text file
- Upload small image (PNG)
- Upload files with special characters in filenames
- Upload files in nested paths
- Upload larger files (1KB)

### Delete Tests (`TestObjectStoreClient_DeleteObject`)
- Delete existing objects
- Delete non-existent objects (idempotent behavior)
- Delete objects in nested paths

### Signed URL Tests (`TestObjectStoreClient_SignedURL`)
- Generate signed URLs for existing objects
- Generate signed URLs for non-existent objects

### Batch Operations (`TestObjectStoreClient_UploadAndDeleteMultiple`)
- Upload multiple objects in sequence
- Verify all uploaded objects exist
- Clean up all uploaded objects

### Content Type Tests (`TestObjectStoreClient_PutObjectWithDifferentExtensions`)
- Test different file extensions (.jpg, .png, .txt, .json, .pdf)
- Verify content-type detection works correctly

### Edge Cases
- **Empty files** (`TestObjectStoreClient_EmptyFile`)
- **Large files** (`TestObjectStoreClient_LargeFile`) - 5MB upload (skipped in short mode)
- **Delete idempotency** (`TestObjectStoreClient_DeleteIdempotency`) - Verify deleting same object multiple times doesn't error

## Requirements

1. **Valid GCS credentials**: The tests require valid Google Cloud Storage credentials
   - File: `sitegrupysanca-firebase-adminsdk-fbsvc-ff7567bd6e.json` (in project root or configs directory)
   - Or application default credentials when running on GCP

2. **GCS bucket**: Access to bucket `sitegrupysanca.firebasestorage.app`

3. **Permissions**: Service account must have:
   - `storage.objects.create`
   - `storage.objects.delete`
   - `storage.objects.get`
   - `storage.buckets.get`

## Test Pattern

All tests follow the same pattern as Firestore tests:

1. **Setup** - Initialize GCS gateway and client
2. **Execute** - Perform operation (upload, delete, etc.)
3. **Cleanup** - Always delete uploaded objects using `defer`
4. **Validate** - Assert expected behavior

Example:
```go
func TestExample(t *testing.T) {
    objectStore, cleanup := setupTestObjectStore(t)
    defer cleanup()

    // Upload object
    url, err := objectStore.PutObject(ctx, "test.txt", data)
    require.NoError(t, err)

    // Cleanup uploaded object
    defer objectStore.DeleteObject(ctx, "test.txt")

    // Validate
    assert.NotEmpty(t, url)
}
```

## Configuration

Tests use the development configuration (`configs/development.yaml`):

```yaml
gcs:
  bucket_name: sitegrupysanca.firebasestorage.app
  credentials_path: sitegrupysanca-firebase-adminsdk-fbsvc-ff7567bd6e.json
  project_id: sitegrupysanca
  base_path: test/images  # All test files go under test/images/
```

## Cleanup

**Important**: All tests clean up after themselves! Objects uploaded during tests are automatically deleted via `defer` statements.

If a test fails or is interrupted, some test objects may remain in the bucket under `test/images/`. You can manually clean them up:

```bash
# Using gcloud CLI
gcloud storage rm -r gs://sitegrupysanca.firebasestorage.app/test/images/
```

## Performance

Typical test run times:
- Individual tests: 0.4s - 2.5s
- Full test suite: ~20-25s
- Large file test: ~3s (skipped with `-short`)

The tests make real API calls to GCS, so network latency affects run time.

## Troubleshooting

### Tests fail with "credentials file not found"
- Ensure credentials file exists in project root or configs directory
- Check `configs/development.yaml` has correct `credentials_path`

### Tests fail with "bucket access denied"
- Verify service account has proper IAM permissions
- Check bucket name is correct in config

### Tests timeout
- Increase timeout with `-timeout` flag: `go test -timeout 5m ...`
- Check network connectivity to Google Cloud Storage

## Related Files

- `internal/gateway/gcs/storage.go` - GCS gateway implementation
- `internal/clients/objectstore.go` - Object store client wrapper
- `configs/config.go` - Configuration handling
