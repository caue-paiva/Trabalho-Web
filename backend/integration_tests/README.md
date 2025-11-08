# Integration Tests

This directory contains HTTP-based integration tests for the Media CMS backend API. These tests make real HTTP requests to a running server and verify the responses.

OBS: Some tests are flaky, especially when it comes to parsing error status-codes 

## Test Files

- **`common_test.go`** - Shared utilities and helper functions
- **`texts_test.go`** - Tests for `/api/v1/texts` endpoints
- **`images_test.go`** - Tests for `/api/v1/images` endpoints
- **`timeline_test.go`** - Tests for `/api/v1/timelineentries` endpoints
- **`events_test.go`** - Tests for `/api/v1/events` endpoint

## Prerequisites

1. **Server must be running** at `http://localhost:8080`
2. **Database (Firestore)** must be configured and accessible
3. **Object Storage (GCS)** must be configured for image upload tests
4. **Network access** to external Events API for events tests

## Running the Tests

### 1. Start the server

In one terminal:

```bash
go run cmd/server/main.go
```

Wait for the log message:
```
Server listening on port 8080
```

### 2. Run all integration tests

In another terminal:

```bash
go test ./integration_tests -v
```

### 3. Run specific test file

```bash
# Run only texts tests
go test ./integration_tests -v -run TestTexts

# Run only images tests
go test ./integration_tests -v -run TestImages

# Run only timeline tests
go test ./integration_tests -v -run TestTimeline

# Run only events tests
go test ./integration_tests -v -run TestEvents
```

### 4. Run specific test

```bash
# Run a single test
go test ./integration_tests -v -run TestTexts_CreateAndGet

# Run tests matching a pattern
go test ./integration_tests -v -run TestImages_.*
```

## Test Coverage

### Texts Endpoints (`texts_test.go`)

✅ **Create and retrieve**
- POST `/api/v1/texts` - Create text
- GET `/api/v1/texts/{slug}` - Get by slug
- GET `/api/v1/texts/id/{id}` - Get by ID

✅ **Update**
- PUT `/api/v1/texts/{id}` - Update text content

✅ **Delete**
- DELETE `/api/v1/texts/{id}` - Delete text

✅ **List and filter**
- GET `/api/v1/texts` - List all texts
- GET `/api/v1/texts/page/slug/{pageSlug}` - Get texts by page

✅ **Error cases**
- 404 for non-existent resources

### Images Endpoints (`images_test.go`)

✅ **Create and retrieve**
- POST `/api/v1/images` - Upload image with base64 data
- GET `/api/v1/images/{id}` - Get image metadata

✅ **Update metadata**
- PUT `/api/v1/images/{id}` - Update name, text, location

✅ **Delete**
- DELETE `/api/v1/images/{id}` - Delete image (DB + storage)

✅ **Gallery operations**
- GET `/api/v1/images/gallery/{slug}` - Get images by gallery

✅ **Validation**
- Object URL accessibility verification
- Invalid base64 handling
- 404 error cases

### Timeline Endpoints (`timeline_test.go`)

✅ **Create and retrieve**
- POST `/api/v1/timelineentries` - Create timeline entry
- GET `/api/v1/timelineentries/{id}` - Get by ID

✅ **Update**
- PUT `/api/v1/timelineentries/{id}` - Update entry

✅ **Delete**
- DELETE `/api/v1/timelineentries/{id}` - Delete entry

✅ **List**
- GET `/api/v1/timelineentries` - List all entries

✅ **Validation**
- Date format validation
- Chronological ordering tests
- 404 error cases

### Events Endpoint (`events_test.go`)

✅ **Basic retrieval**
- GET `/api/v1/events` - Get all events

✅ **Query parameters**
- `?limit=N` - Limit results
- `?orderBy=startsAt` - Order by field
- `?desc=true` - Descending order
- Combined parameters

✅ **Response structure**
- Verify all expected fields present
- Handle empty results gracefully

## Test Characteristics

### Cleanup Strategy

All tests that create resources **automatically clean up** after themselves using `defer`:

```go
defer func() {
    resp := MakeRequest(t, "DELETE", "/texts/"+created.ID, nil)
    resp.Body.Close()
}()
```

This ensures:
- No test pollution between runs
- Database stays clean
- Object storage doesn't accumulate test files

### Unique Identifiers

Tests use timestamp-based unique slugs to avoid conflicts:

```go
slug := GenerateUniqueSlug("test-prefix")
// Result: "test-prefix-1731096377"
```

### Test Independence

- Each test is independent and can run in isolation
- Tests don't rely on execution order
- Tests create their own test data

## Configuration

### Change Server URL

Edit `common_test.go`:

```go
const BaseURL = "http://localhost:8080/api/v1"
```

### Change Request Timeout

Edit `common_test.go`:

```go
const RequestTimeout = 10 * time.Second
```

## Common Patterns

### Making a Request

```go
resp := MakeRequest(t, "POST", "/texts", requestBody)
AssertStatusCode(t, resp, http.StatusCreated)

var response TextResponse
ParseJSONResponse(t, resp, &response)
```

### Cleanup Pattern

```go
// Create resource
resp := MakeRequest(t, "POST", "/texts", createReq)
var created TextResponse
ParseJSONResponse(t, resp, &created)

// Always cleanup
defer func() {
    resp := MakeRequest(t, "DELETE", "/texts/"+created.ID, nil)
    resp.Body.Close()
}()

// ... test assertions ...
```

## Troubleshooting

### "connection refused" errors

**Problem:** Server is not running
**Solution:** Start the server with `go run cmd/server/main.go`

### Tests timeout

**Problem:** Server is slow to respond or hung
**Solution:**
- Check server logs for errors
- Increase `RequestTimeout` in `common_test.go`
- Verify database and storage are accessible

### "404 Not Found" errors

**Problem:** Wrong base URL or routes changed
**Solution:**
- Verify server is running on port 8080
- Check `BaseURL` in `common_test.go`
- Verify routes match `internal/http/router.go`

### Image tests fail with GCS errors

**Problem:** Object storage not configured
**Solution:**
- Verify GCS credentials in `configs/development.yaml`
- Check service account has proper permissions
- Ensure bucket exists and is accessible

### Event tests return empty arrays

**Problem:** External Events API is down or has no data
**Solution:** This is expected behavior - tests handle gracefully

## Running with Different Environments

### Development (default)

```bash
RUNTIME_ENV=development go test ./integration_tests -v
```

### Production

```bash
RUNTIME_ENV=production BASE_URL=https://api.example.com/api/v1 go test ./integration_tests -v
```

## Continuous Integration

These tests are suitable for CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Start server
  run: go run cmd/server/main.go &

- name: Wait for server
  run: sleep 5

- name: Run integration tests
  run: go test ./integration_tests -v
```

## Performance Notes

- **Average test time**: 0.5-2 seconds per test
- **Full suite**: ~30-60 seconds (depends on network and DB)
- **Slowest tests**: Image upload/download (GCS latency)
- **Fastest tests**: Events endpoint (read-only, cached)

## Best Practices

1. **Always cleanup** - Use `defer` to delete created resources
2. **Use unique identifiers** - Prevent test conflicts
3. **Assert on critical fields** - Don't assert on timestamps (flaky)
4. **Handle flakiness** - External APIs may be unavailable
5. **Keep tests independent** - No shared state between tests

## Future Enhancements

Possible additions:
- [ ] Authentication/authorization tests (when implemented)
- [ ] Pagination tests (limit, offset)
- [ ] Concurrent request tests (race conditions)
- [ ] Performance/load tests
- [ ] WebSocket tests (if added)
- [ ] File upload size limit tests
- [ ] Rate limiting tests (if implemented)

## Contributing

When adding new tests:

1. Follow existing patterns in `*_test.go` files
2. Add cleanup with `defer`
3. Use `GenerateUniqueSlug()` for unique identifiers
4. Update this README with new test coverage
5. Ensure tests are independent and idempotent
