# Configuration Service

A flexible YAML-based configuration service that supports environment-specific configs and struct unmarshaling.

## Features

- ✅ Environment-based config loading (development/production)
- ✅ Nested key access with dot notation
- ✅ Struct unmarshaling with yaml tags
- ✅ Built-in Go stdlib (uses `gopkg.in/yaml.v3`)

## Environment Selection

The config service reads from different YAML files based on the `RUNTIME_ENV` environment variable:

- **`RUNTIME_ENV=development`** → reads `configs/development.yaml` (default)
- **`RUNTIME_ENV=production`** → reads `configs/production.yaml`

## Usage Examples

### 1. Initialize Config Service

```go
import "backend/configs"

func main() {
    // Reads development.yaml by default (or production.yaml if RUNTIME_ENV=production)
    config, err := configs.NewConfigService()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Use config...
}
```

### 2. Get Simple Values by Key

```go
// Get a top-level value
firebasePath, err := config.GetConfig("firebase_config_path")
if err != nil {
    log.Fatal(err)
}
fmt.Println(firebasePath) // "sitegrupysanca-firebase-adminsdk-fbsvc-ff7567bd6e.json"
```

### 3. Get Nested Values with Dot Notation

```go
// development.yaml:
// collections:
//   texts: test_texts
//   images: test_images

textsCollection, err := config.GetConfig("collections.texts")
if err != nil {
    log.Fatal(err)
}
fmt.Println(textsCollection) // "test_texts"

imagesCollection, err := config.GetConfig("collections.images")
fmt.Println(imagesCollection) // "test_images"
```

### 4. Unmarshal Config Section into Struct

```go
// Define struct with yaml tags
type Collections struct {
    Texts     string `yaml:"texts"`
    Images    string `yaml:"images"`
    Timelines string `yaml:"timelines"`
}

// Unmarshal "collections" section
var cols Collections
err := config.UnmarshalKey("collections", &cols)
if err != nil {
    log.Fatal(err)
}

fmt.Println(cols.Texts)     // "test_texts"
fmt.Println(cols.Images)    // "test_images"
fmt.Println(cols.Timelines) // "test_timelines"
```

### 5. Complete Example with Firestore Integration

```go
package main

import (
    "context"
    "log"

    "backend/configs"
    "backend/internal/repository/firestore"
)

func main() {
    ctx := context.Background()

    // Load configuration
    config, err := configs.NewConfigService()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Get Firebase credentials path
    credPath, err := config.GetConfig("firebase_config_path")
    if err != nil {
        log.Fatalf("Failed to get firebase_config_path: %v", err)
    }

    // Unmarshal collection names
    type Collections struct {
        Texts     string `yaml:"texts"`
        Images    string `yaml:"images"`
        Timelines string `yaml:"timelines"`
    }

    var cols Collections
    if err := config.UnmarshalKey("collections", &cols); err != nil {
        log.Fatalf("Failed to unmarshal collections: %v", err)
    }

    // Create Firestore collection names
    collections := firestore.CollectionNames{
        Texts:           cols.Texts,
        Images:          cols.Images,
        TimelineEntries: cols.Timelines,
    }

    // Initialize Firestore
    firestoreConfig := firestore.FirestoreConfig{
        ProjectID:       "your-project-id",
        CredentialsPath: credPath.(string),
        Collections:     collections,
    }

    client, err := firestore.NewFirestoreClient(ctx, firestoreConfig)
    if err != nil {
        log.Fatalf("Failed to create Firestore client: %v", err)
    }
    defer client.Close()

    // Create DB repository
    dbRepo := firestore.NewDBRepository(client, collections)

    // Use dbRepo...
}
```

## Config File Format

### development.yaml

```yaml
# Config for local development and testing

# Firebase collection paths/names
collections:
  texts: test_texts
  timelines: test_timelines
  images: test_images

firebase_config_path: sitegrupysanca-firebase-adminsdk-fbsvc-ff7567bd6e.json
```

### production.yaml

```yaml
# Production config

collections:
  texts: texts
  timelines: timeline_entries
  images: images

firebase_config_path: /secrets/firebase-credentials.json
```

## API Reference

### `NewConfigService() (ConfigClient, error)`

Creates a new config service. Reads the appropriate YAML file based on `RUNTIME_ENV`.

**Returns:**
- `ConfigClient` - Config service instance
- `error` - Error if config file cannot be read or parsed

### `GetConfig(cfgName string) (any, error)`

Gets a config value by key. Supports nested keys with dot notation.

**Parameters:**
- `cfgName` - Config key (e.g., `"firebase_config_path"` or `"collections.texts"`)

**Returns:**
- `any` - Config value (can be string, map, slice, etc.)
- `error` - Error if key not found

**Examples:**
```go
val, _ := config.GetConfig("firebase_config_path")         // Simple key
val, _ := config.GetConfig("collections.texts")            // Nested key
val, _ := config.GetConfig("database.host")                // Deep nesting
```

### `UnmarshalKey(key string, target any) error`

Unmarshals a config section into a struct using yaml tags.

**Parameters:**
- `key` - Config section key (e.g., `"collections"`)
- `target` - Pointer to struct to fill

**Returns:**
- `error` - Error if key not found or unmarshaling fails

**Example:**
```go
type DBConfig struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
}

var db DBConfig
err := config.UnmarshalKey("database", &db)
```

## Error Handling

All methods return descriptive errors:

```go
val, err := config.GetConfig("non.existent.key")
if err != nil {
    // Error: "config key 'non.existent.key' not found"
    log.Fatal(err)
}

err = config.UnmarshalKey("collections", nil)
if err != nil {
    // Error: "target cannot be nil"
    log.Fatal(err)
}
```

## Testing

Run tests with:

```bash
go test ./configs -v
```

All tests use the `development.yaml` file and verify:
- ✅ Config loading from different environments
- ✅ Simple key retrieval
- ✅ Nested key access with dot notation
- ✅ Struct unmarshaling
- ✅ Error handling
