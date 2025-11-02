package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ConfigClient provides access to configuration values
type ConfigClient interface {
	// GetConfig returns a config value by key (supports nested keys with dots, e.g., "collections.texts")
	GetConfig(cfgName string) (any, error)

	// UnmarshalKey unmarshals a config section into a struct pointer using yaml tags
	// Example: UnmarshalKey("collections", &collectionsStruct)
	UnmarshalKey(key string, target any) error
}

type configService struct {
	data map[string]any
	env  string
}

// NewConfigService creates a new config service
// Reads config from YAML files based on RUNTIME_ENV (defaults to development)
func NewConfigService() (ConfigClient, error) {
	// Get runtime environment (development or production)
	env := os.Getenv("RUNTIME_ENV")
	if env == "" {
		env = "development" // default
	}

	// Validate environment
	if env != "development" && env != "production" {
		env = "development"
	}

	// Determine config file path
	configFile := fmt.Sprintf("%s.yaml", env)

	// Try multiple locations (for different working directories)
	possiblePaths := []string{
		filepath.Join("configs", configFile),  // From project root
		configFile,                            // From configs directory (for tests)
		filepath.Join("..", "configs", configFile), // One level up
	}

	var data []byte
	var configPath string
	for _, path := range possiblePaths {
		var err error
		data, err = os.ReadFile(path)
		if err == nil {
			configPath = path
			break
		}
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("failed to read config file %s from any location", configFile)
	}

	// Parse YAML
	var configData map[string]any
	if err := yaml.Unmarshal(data, &configData); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	return &configService{
		data: configData,
		env:  env,
	}, nil
}

// GetConfig returns a config value by key path
// Supports nested keys using dot notation (e.g., "collections.texts")
func (s *configService) GetConfig(cfgName string) (any, error) {
	if cfgName == "" {
		return nil, fmt.Errorf("config name cannot be empty")
	}

	// Split key path by dots
	keys := strings.Split(cfgName, ".")

	// Navigate through nested map
	var current any = s.data

	for i, key := range keys {
		// Assert current value is a map
		currentMap, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("config key '%s' is not a map at level %d", strings.Join(keys[:i], "."), i)
		}

		// Get value for this key
		value, exists := currentMap[key]
		if !exists {
			return nil, fmt.Errorf("config key '%s' not found", cfgName)
		}

		current = value
	}

	return current, nil
}

// UnmarshalKey unmarshals a specific config section into a struct
// The target must be a pointer to a struct with yaml tags
//
// Example:
//   type Collections struct {
//       Texts string `yaml:"texts"`
//       Images string `yaml:"images"`
//   }
//   var cols Collections
//   err := config.UnmarshalKey("collections", &cols)
func (s *configService) UnmarshalKey(key string, target any) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	// Get the config section
	value, err := s.GetConfig(key)
	if err != nil {
		return err
	}

	// Marshal the value back to YAML bytes (so we can unmarshal into struct)
	yamlBytes, err := yaml.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal config section '%s': %w", key, err)
	}

	// Unmarshal into target struct
	if err := yaml.Unmarshal(yamlBytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal config section '%s' into target: %w", key, err)
	}

	return nil
}
