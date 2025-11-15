package integration_tests

import (
	"backend/configs"
	"backend/internal/http/mapper"
	"backend/internal/platform/auth"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthOptional_PostTextWithoutAuth(t *testing.T) {
	// 1. Call config function GetAuthLevel and assert it's optional
	config, err := configs.NewConfigService()
	require.NoError(t, err, "Failed to initialize config service")

	authLevel := config.GetAuthLevel()
	assert.Equal(t, auth.AuthOptional, authLevel, "Auth level should be optional")

	// 2. Make POST request to text endpoint without an auth token
	slug := GenerateUniqueSlug("auth-optional-test")
	createReq := mapper.CreateTextRequest{
		Slug:    slug,
		Content: "Test content for auth optional test",
	}

	resp := MakeRequest(t, "POST", "/texts", createReq)

	// 3. Assert success response
	AssertStatusCode(t, resp, http.StatusCreated)

	var created mapper.TextResponse
	ParseJSONResponse(t, resp, &created)

	// Validate created text
	assert.NotEmpty(t, created.ID, "Text should have an ID")
	assert.Equal(t, slug, created.Slug)
	assert.Equal(t, createReq.Content, created.Content)

	// 4. Clean up
	defer func() {
		resp := MakeRequest(t, "DELETE", "/texts/"+created.ID, nil)
		resp.Body.Close()
	}()
}