package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewConfigService_Development tests loading development config
func TestNewConfigService_Development(t *testing.T) {
	// Ensure RUNTIME_ENV is not set (should default to development)
	os.Unsetenv("RUNTIME_ENV")

	config, err := NewConfigService()
	require.NoError(t, err, "Should successfully load development config")
	require.NotNil(t, config)
}

// TestNewConfigService_Production tests loading production config
func TestNewConfigService_Production(t *testing.T) {
	// Set RUNTIME_ENV to production
	os.Setenv("RUNTIME_ENV", "production")
	defer os.Unsetenv("RUNTIME_ENV")

	config, err := NewConfigService()
	require.NoError(t, err, "Should successfully load production config")
	require.NotNil(t, config)
}

// TestGetConfig_SimpleKey tests getting a simple nested config value
func TestGetConfig_SimpleKey(t *testing.T) {
	os.Unsetenv("RUNTIME_ENV") // Use development config

	config, err := NewConfigService()
	require.NoError(t, err)

	// Get firebase.project_id from development.yaml
	value, err := config.GetConfig("firebase.project_id")
	require.NoError(t, err)
	assert.Equal(t, "sitegrupysanca", value)
}

// TestGetConfig_NestedKey tests getting nested config values with dot notation
func TestGetConfig_NestedKey(t *testing.T) {
	os.Unsetenv("RUNTIME_ENV")

	config, err := NewConfigService()
	require.NoError(t, err)

	// Get collections.texts from development.yaml
	value, err := config.GetConfig("collections.texts")
	require.NoError(t, err)
	assert.Equal(t, "test_texts", value)

	// Get collections.images
	value, err = config.GetConfig("collections.images")
	require.NoError(t, err)
	assert.Equal(t, "test_images", value)

	// Get collections.timelines
	value, err = config.GetConfig("collections.timelines")
	require.NoError(t, err)
	assert.Equal(t, "test_timelines", value)
}

// TestGetConfig_NonExistentKey tests error handling for missing keys
func TestGetConfig_NonExistentKey(t *testing.T) {
	os.Unsetenv("RUNTIME_ENV")

	config, err := NewConfigService()
	require.NoError(t, err)

	_, err = config.GetConfig("non.existent.key")
	assert.Error(t, err, "Should return error for non-existent key")
	assert.Contains(t, err.Error(), "not found")
}

// TestUnmarshalKey_CollectionsStruct tests unmarshaling config section into struct
func TestUnmarshalKey_CollectionsStruct(t *testing.T) {
	os.Unsetenv("RUNTIME_ENV")

	config, err := NewConfigService()
	require.NoError(t, err)

	// Define struct matching the collections config
	type Collections struct {
		Texts     string `yaml:"texts"`
		Images    string `yaml:"images"`
		Timelines string `yaml:"timelines"`
	}

	var cols Collections
	err = config.UnmarshalKey("collections", &cols)
	require.NoError(t, err)

	// Verify values were correctly unmarshaled
	assert.Equal(t, "test_texts", cols.Texts)
	assert.Equal(t, "test_images", cols.Images)
	assert.Equal(t, "test_timelines", cols.Timelines)
}

// TestUnmarshalKey_EmptyKey tests error handling for empty key
func TestUnmarshalKey_EmptyKey(t *testing.T) {
	os.Unsetenv("RUNTIME_ENV")

	config, err := NewConfigService()
	require.NoError(t, err)

	type Collections struct {
		Texts     string `yaml:"texts"`
		Images    string `yaml:"images"`
		Timelines string `yaml:"timelines"`
	}

	var fc Collections
	err = config.UnmarshalKey("", &fc)
	// This should fail because empty key is not allowed
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

// TestUnmarshalKey_NilTarget tests error handling for nil target
func TestUnmarshalKey_NilTarget(t *testing.T) {
	os.Unsetenv("RUNTIME_ENV")

	config, err := NewConfigService()
	require.NoError(t, err)

	err = config.UnmarshalKey("collections", nil)
	assert.Error(t, err, "Should return error for nil target")
	assert.Contains(t, err.Error(), "cannot be nil")
}
