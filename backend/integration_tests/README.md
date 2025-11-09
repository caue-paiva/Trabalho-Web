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

## CURL Examples

This section provides example CURL commands to manually test all API endpoints. The base URL is `http://localhost:8080/api/v1` (adjust if your server runs on a different port).

### Texts Endpoints

#### List All Texts
```bash
curl -X GET http://localhost:8080/api/v1/texts
```

#### Get Text by Slug
```bash
curl -X GET http://localhost:8080/api/v1/texts/my-text-slug
```

#### Get Text by ID
```bash
curl -X GET http://localhost:8080/api/v1/texts/id/abc123def456
```

#### Get Texts by Page Slug
```bash
curl -X GET http://localhost:8080/api/v1/texts/page/slug/historia
```

#### Get Texts by Page ID
```bash
curl -X GET http://localhost:8080/api/v1/texts/page/page-id-123
```

#### Create Text
```bash
curl -X POST http://localhost:8080/api/v1/texts \
  -H "Content-Type: application/json" \
  -d '{
    "slug": "my-text-slug",
    "content": "This is the text content",
    "page_slug": "historia"
  }'
```

#### Update Text
```bash
curl -X PUT http://localhost:8080/api/v1/texts/abc123def456 \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Updated text content",
    "page_slug": "updated-page"
  }'
```

#### Delete Text
```bash
curl -X DELETE http://localhost:8080/api/v1/texts/abc123def456
```

### Images Endpoints

#### Get Image by ID
```bash
curl -X GET http://localhost:8080/api/v1/images/abc123def456
```

#### Get Images by Gallery Slug
```bash
curl -X GET http://localhost:8080/api/v1/images/gallery/my-gallery-slug
```

#### Create Image (Upload)
```bash
# First, encode your image to base64:
# On macOS/Linux:
IMAGE_BASE64=$(base64 -i path/to/image.png)

# Then create the image:
curl -X POST http://localhost:8080/api/v1/images \
  -H "Content-Type: application/json" \
  -d "{
    \"slug\": \"my-gallery-slug\",
    \"name\": \"My Image Name\",
    \"text\": \"Image description\",
    \"date\": \"2024-01-15\",
    \"location\": \"São Carlos, SP\",
    \"data\": \"$IMAGE_BASE64\"
  }"
```

**Note:** For testing, you can use a tiny 1x1 pixel PNG in base64:
```bash
curl -X POST http://localhost:8080/api/v1/images \
  -H "Content-Type: application/json" \
  -d '{
    "slug": "test-gallery",
    "name": "Test Image",
    "text": "Test description",
    "data": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8DwHwAFBQIAX8jx0gAAAABJRU5ErkJggg=="
  }'
```

#### Update Image Metadata (without new image)
```bash
curl -X PUT http://localhost:8080/api/v1/images/abc123def456 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Image Name",
    "text": "Updated description",
    "location": "Updated Location"
  }'
```

#### Update Image (with new image data)
```bash
# Encode new image to base64 first
IMAGE_BASE64=$(base64 -i path/to/new-image.png)

curl -X PUT http://localhost:8080/api/v1/images/abc123def456 \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Updated Image\",
    \"data\": \"$IMAGE_BASE64\"
  }"
```

#### Delete Image
```bash
curl -X DELETE http://localhost:8080/api/v1/images/abc123def456
```

### Timeline Endpoints

#### List All Timeline Entries
```bash
curl -X GET http://localhost:8080/api/v1/timelineentries
```

#### Get Timeline Entry by ID
```bash
curl -X GET http://localhost:8080/api/v1/timelineentries/abc123def456
```

#### Create Timeline Entry
```bash
curl -X POST http://localhost:8080/api/v1/timelineentries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "GrupySanca Meetup",
    "text": "An important milestone in our history",
    "location": "São Carlos, SP",
    "date": "2024-11-08T12:00:00Z"
  }'
```

**Note:** Date must be in RFC3339 format (ISO 8601). Examples:
- `2024-11-08T12:00:00Z` (UTC)
- `2024-11-08T12:00:00-03:00` (with timezone)

#### Update Timeline Entry
```bash
curl -X PUT http://localhost:8080/api/v1/timelineentries/abc123def456 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Event Name",
    "text": "Updated description",
    "location": "Updated Location",
    "date": "2024-12-20T14:00:00Z"
  }'
```

#### Delete Timeline Entry
```bash
curl -X DELETE http://localhost:8080/api/v1/timelineentries/abc123def456
```

### Events Endpoints

#### Get All Events
```bash
curl -X GET http://localhost:8080/api/v1/events
```

#### Get Events with Limit
```bash
curl -X GET "http://localhost:8080/api/v1/events?limit=5"
```

#### Get Events with Limit (larger)
```bash
curl -X GET "http://localhost:8080/api/v1/events?limit=100"
```

**Note:** The events endpoint fetches data from an external API. Query parameters like `limit` are supported, but `orderBy` and `desc` may not be fully supported by the external API.

### Error Handling Examples

#### Test 404 Not Found (Text)
```bash
curl -X GET http://localhost:8080/api/v1/texts/non-existent-slug-12345
# Expected: 404 Not Found
```

#### Test 404 Not Found (Image)
```bash
curl -X GET http://localhost:8080/api/v1/images/non-existent-id-12345
# Expected: 404 Not Found
```

#### Test 404 Not Found (Timeline)
```bash
curl -X GET http://localhost:8080/api/v1/timelineentries/non-existent-id-12345
# Expected: 404 Not Found
```

#### Test Invalid Base64 (Image)
```bash
curl -X POST http://localhost:8080/api/v1/images \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Invalid Image",
    "text": "This has invalid base64 data",
    "data": "not-valid-base64!@#$%"
  }'
# Expected: 400 Bad Request
```

#### Test Invalid Date Format (Timeline)
```bash
curl -X POST http://localhost:8080/api/v1/timelineentries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Invalid Date Event",
    "text": "This has an invalid date",
    "date": "not-a-valid-date"
  }'
# Expected: 400 Bad Request
```

### Pretty-Printing JSON Responses

For better readability, pipe responses through `jq` (if installed):

```bash
curl -X GET http://localhost:8080/api/v1/texts | jq
```

Or use Python:
```bash
curl -X GET http://localhost:8080/api/v1/texts | python -m json.tool
```

### Saving Responses to Files

```bash
# Save response to file
curl -X GET http://localhost:8080/api/v1/texts -o response.json

# Save with pretty formatting
curl -X GET http://localhost:8080/api/v1/texts | jq > response.json
```

### Testing with Different Base URLs

If your server runs on a different host/port, set a variable:

```bash
BASE_URL="http://localhost:8080/api/v1"
curl -X GET "$BASE_URL/texts"
```

For production or staging:
```bash
BASE_URL="https://api.example.com/api/v1"
curl -X GET "$BASE_URL/texts"
```

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
