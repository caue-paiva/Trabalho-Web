# Firestore Repository

This package provides a Firestore implementation of the `service.DBPort` interface.

## Setup

### 1. Create a Firebase Project

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Create a new project or select an existing one
3. Enable Firestore Database:
   - Go to Firestore Database section
   - Click "Create database"
   - Choose "Production mode" or "Test mode" (for development)
   - Select a region close to your users

### 2. Get Service Account Credentials

1. In Firebase Console, go to **Project Settings** > **Service Accounts**
2. Click "Generate new private key"
3. Save the JSON file securely (e.g., `configs/serviceAccountKey.json`)
4. **IMPORTANT**: Add this file to `.gitignore` - never commit credentials!

### 3. Environment Setup

You can authenticate in two ways:

#### Option A: Service Account File (Recommended for Development)
```go
client, err := firestore.NewFirestoreClient(
    ctx,
    "your-project-id",
    "path/to/serviceAccountKey.json",
)
```

#### Option B: Application Default Credentials (Recommended for Production/GCP)
```go
// Set GOOGLE_APPLICATION_CREDENTIALS environment variable
// export GOOGLE_APPLICATION_CREDENTIALS="/path/to/serviceAccountKey.json"

client, err := firestore.NewFirestoreClient(
    ctx,
    "your-project-id",
    "", // Empty string uses default credentials
)
```

### 4. Usage Example

```go
package main

import (
    "context"
    "log"

    "backend/internal/repository/firestore"
    "backend/internal/service"
)

func main() {
    ctx := context.Background()

    // Initialize Firestore client
    firestoreClient, err := firestore.NewFirestoreClient(
        ctx,
        "your-project-id",
        "configs/serviceAccountKey.json",
    )
    if err != nil {
        log.Fatalf("Failed to create Firestore client: %v", err)
    }
    defer firestoreClient.Close()

    // Create DB repository
    dbRepo := firestore.NewDBRepository(firestoreClient)

    // Create services
    textContentService := service.NewTextContentService(dbRepo)
    imageService := service.NewImageService(dbRepo, objectStorePort)

    // Use services...
}
```

## Collections Structure

The repository creates the following Firestore collections:

- **texts**: Text content blocks
  - Fields: `id`, `slug`, `content`, `pageID`, `pageSlug`, `createdAt`, `updatedAt`, `lastUpdatedBy`

- **images**: Image metadata
  - Fields: `id`, `slug`, `objectURL`, `name`, `text`, `date`, `location`, `createdAt`, `updatedAt`, `lastUpdatedBy`

- **timeline_entries**: Timeline events
  - Fields: `id`, `name`, `text`, `location`, `date`, `createdAt`, `updatedAt`, `lastUpdatedBy`

## Security Rules (Production)

For production, set up Firestore security rules:

```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    // Allow public read access to all collections
    match /{collection}/{document} {
      allow read: if true;
      allow write: if request.auth != null; // Require authentication for writes
    }
  }
}
```

## Local Development with Emulator (Optional)

For local development without using production Firestore:

```bash
# Install Firebase emulator
npm install -g firebase-tools

# Initialize Firebase in your project
firebase init emulators

# Start emulator
firebase emulators:start
```

Then connect to the emulator:
```go
// Set environment variable
os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")

client, err := firestore.NewFirestoreClient(ctx, "demo-project", "")
```

## Error Handling

The repository returns descriptive errors:
- "not found" errors when documents don't exist
- Wrapped errors with context for all operations
- Iterator errors are properly handled

## Best Practices

1. **Always close the Firestore client** when done:
   ```go
   defer firestoreClient.Close()
   ```

2. **Use context for timeouts**:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
   defer cancel()
   ```

3. **Handle errors appropriately** - check for "not found" vs other errors

4. **Don't commit service account keys** - use environment variables or secret management
