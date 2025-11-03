package firestore

import (
	"context"

	"backend/configs"
	"cloud.google.com/go/firestore"
)

// FirebaseConfigProvider defines the interface for accessing Firebase configuration
type FirebaseConfigProvider interface {
	// GetFirebaseConfig returns the Firebase configuration
	GetFirebaseConfig() (configs.FirebaseConfig, error)

	// GetCollections returns the Firestore collection names
	GetCollections() (configs.Collections, error)

	// GetCredentialsJSON reads the Firebase credentials JSON file and returns its bytes
	GetCredentialsJSON(filename string) ([]byte, error)
}

// NewFirestoreClientWithProvider creates a new Firestore client using a config provider
func NewFirestoreClientWithProvider(ctx context.Context, provider FirebaseConfigProvider) (*firestore.Client, error) {
	// Get Firebase config from provider
	fbConfig, err := provider.GetFirebaseConfig()
	if err != nil {
		return nil, err
	}

	// Read credentials
	credentialsJSON, err := provider.GetCredentialsJSON(fbConfig.CredentialsPath)
	if err != nil {
		return nil, err
	}

	// Get collections
	collections, err := provider.GetCollections()
	if err != nil {
		return nil, err
	}

	// Create Firestore config
	firestoreConfig := FirestoreConfig{
		ProjectID:       fbConfig.ProjectID,
		CredentialsJSON: credentialsJSON,
		Collections: CollectionNames{
			Texts:           collections.Texts,
			Images:          collections.Images,
			TimelineEntries: collections.Timelines,
		},
	}

	// Create and return Firestore client
	return NewFirestoreClient(ctx, firestoreConfig)
}

// NewDBRepositoryWithProvider creates a new DBRepository using a config provider
func NewDBRepositoryWithProvider(ctx context.Context, provider FirebaseConfigProvider) (*DBRepository, error) {
	// Create Firestore client using the provider
	client, err := NewFirestoreClientWithProvider(ctx, provider)
	if err != nil {
		return nil, err
	}

	// Get collections
	collections, err := provider.GetCollections()
	if err != nil {
		return nil, err
	}

	// Convert Collections to CollectionNames
	collectionNames := CollectionNames{
		Texts:           collections.Texts,
		Images:          collections.Images,
		TimelineEntries: collections.Timelines,
	}

	// Create and return DB repository
	return NewDBRepository(client, collectionNames), nil
}
