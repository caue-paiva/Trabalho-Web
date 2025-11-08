package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	// BaseURL is the base URL of the API server
	// Change this if your server runs on a different port
	BaseURL = "http://localhost:8080/api/v1"

	// Timeout for HTTP requests
	RequestTimeout = 10 * time.Second
)

// HTTPClient is a configured HTTP client for integration tests
var HTTPClient = &http.Client{
	Timeout: RequestTimeout,
}

// MakeRequest makes an HTTP request and returns the response
func MakeRequest(t *testing.T, method, path string, body interface{}) *http.Response {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		require.NoError(t, err, "Failed to marshal request body")
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := BaseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	require.NoError(t, err, "Failed to create request")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := HTTPClient.Do(req)
	require.NoError(t, err, "Failed to make request to %s %s", method, path)

	return resp
}

// ParseJSONResponse parses a JSON response into the target struct
func ParseJSONResponse(t *testing.T, resp *http.Response, target interface{}) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	err = json.Unmarshal(body, target)
	require.NoError(t, err, "Failed to unmarshal JSON response: %s", string(body))
}

// AssertStatusCode asserts that the response has the expected status code
func AssertStatusCode(t *testing.T, resp *http.Response, expectedStatus int) {
	if resp.StatusCode != expectedStatus {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status %d, got %d. Response body: %s",
			expectedStatus, resp.StatusCode, string(body))
	}
}

// GenerateUniqueSlug generates a unique slug with timestamp
func GenerateUniqueSlug(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().Unix())
}

// WaitForServer waits for the server to be ready
func WaitForServer(t *testing.T, maxRetries int) {
	for i := 0; i < maxRetries; i++ {
		resp, err := HTTPClient.Get(BaseURL + "/texts")
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatal("Server not ready after maximum retries")
}
