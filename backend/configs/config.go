package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FirebaseConfig holds Firebase-specific configuration loaded from YAML
type FirebaseConfig struct {
	ProjectID       string `yaml:"project_id"`
	CredentialsPath string `yaml:"credentials_path"`
}

// Collections holds the names of Firestore collections loaded from YAML
type Collections struct {
	Texts     string `yaml:"texts"`
	Images    string `yaml:"images"`
	Timelines string `yaml:"timelines"`
}

// GCSConfig holds Google Cloud Storage configuration
type GCSConfig struct {
	BucketName             string `yaml:"bucket_name"`
	CredentialsPath        string `yaml:"credentials_path"`
	CredentialsJSON        []byte `yaml:"-"` // Populated by GetCredentialsJSON, not from YAML
	ProjectID              string `yaml:"project_id"`
	MakePublic             bool   `yaml:"make_public"`
	SignedURLExpiryMinutes int    `yaml:"signed_url_expiry_minutes"`
	BasePath               string `yaml:"base_path"` // Base path within bucket for all objects (e.g., "images", "media/uploads")
}

// ConfigClient provides access to configuration values
type ConfigClient interface {
	// GetConfig returns a config value by key (supports nested keys with dots, e.g., "collections.texts")
	GetConfig(cfgName string) (any, error)

	// UnmarshalKey unmarshals a config section into a struct pointer using yaml tags
	// Example: UnmarshalKey("collections", &collectionsStruct)
	UnmarshalKey(key string, target any) error

	// GetCredentialsJSON reads the Firebase credentials JSON file and returns its bytes
	// It looks for the file in the configs directory first, then in the project root
	GetCredentialsJSON(filename string) ([]byte, error)

	// GetFirebaseConfig returns the Firebase configuration
	GetFirebaseConfig() (FirebaseConfig, error)

	// GetCollections returns the Firestore collection names
	GetCollections() (Collections, error)

	// GetGCSConfig returns the Google Cloud Storage configuration
	GetGCSConfig() (GCSConfig, error)
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
		filepath.Join("configs", configFile),             // From project root
		configFile,                                       // From configs directory
		filepath.Join("..", "configs", configFile),       // One level up
		filepath.Join("../..", "configs", configFile),    // Two levels up
		filepath.Join("../../..", "configs", configFile), // Three levels up (for deep tests like firestore)
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
//
//	type Collections struct {
//	    Texts string `yaml:"texts"`
//	    Images string `yaml:"images"`
//	}
//	var cols Collections
//	err := config.UnmarshalKey("collections", &cols)
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

// GetCredentialsJSON reads the Firebase credentials JSON file and returns its bytes
// It looks for the file in the configs directory first, then in the project root
func (s *configService) GetCredentialsJSON(filename string) ([]byte, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename cannot be empty")
	}

	// If it's an absolute path, read it directly
	if filepath.IsAbs(filename) {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read credentials file %s: %w", filename, err)
		}
		return data, nil
	}

	// Find project root by walking up to find go.mod
	// If go.mod not found (e.g., in Docker container), use current working directory
	projectRoot, err := findProjectRoot()
	if err != nil {
		// Fallback to current working directory (for Docker containers without go.mod)
		projectRoot, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to find project root or get working directory: %w", err)
		}
	}

	// Try configs directory first, then project root
	possiblePaths := []string{
		filepath.Join(projectRoot, "configs", filename),
		filepath.Join(projectRoot, filename),
		// Also try relative to current working directory (for Docker)
		filepath.Join("configs", filename),
		filename, // Direct filename (if in same directory)
	}

	for _, path := range possiblePaths {
		data, err := os.ReadFile(path)
		if err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("credentials file %s not found in configs directory or project root", filename)
}

// findProjectRoot walks up the directory tree to find go.mod
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check if go.mod exists in current directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding go.mod
			return "", fmt.Errorf("could not find project root (go.mod not found)")
		}
		dir = parent
	}
}

// GetFirebaseConfig returns the Firebase configuration
func (s *configService) GetFirebaseConfig() (FirebaseConfig, error) {
	var config FirebaseConfig
	if err := s.UnmarshalKey("firebase", &config); err != nil {
		return FirebaseConfig{}, err
	}
	return config, nil
}

// GetCollections returns the Firestore collection names
func (s *configService) GetCollections() (Collections, error) {
	var collections Collections
	if err := s.UnmarshalKey("collections", &collections); err != nil {
		return Collections{}, err
	}
	return collections, nil
}

// GetGCSConfig returns the Google Cloud Storage configuration
// If CredentialsPath is specified, it will also populate CredentialsJSON with the file contents
func (s *configService) GetGCSConfig() (GCSConfig, error) {
	var config GCSConfig
	if err := s.UnmarshalKey("gcs", &config); err != nil {
		return GCSConfig{}, err
	}

	// If credentials path is specified, read the credentials file
	if config.CredentialsPath != "" {
		credentialsJSON, err := s.GetCredentialsJSON(config.CredentialsPath)
		if err != nil {
			return GCSConfig{}, fmt.Errorf("failed to read credentials file: %w", err)
		}
		config.CredentialsJSON = credentialsJSON
	}

	return config, nil
}
