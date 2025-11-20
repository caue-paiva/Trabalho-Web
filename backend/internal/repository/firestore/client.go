package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// CollectionNames holds the names of Firestore collections
type CollectionNames struct {
	Texts           string
	Images          string
	TimelineEntries string
	GaleryEvents    string
}

// FirestoreConfig holds configuration for Firestore client initialization
type FirestoreConfig struct {
	ProjectID       string
	CredentialsJSON []byte
	Collections     CollectionNames
}

// NewFirestoreClient creates a new Firestore client
// config.CredentialsJSON should contain the Firebase service account JSON key bytes
func NewFirestoreClient(ctx context.Context, config FirestoreConfig) (*firestore.Client, error) {
	projectID := config.ProjectID
	credentialsJSON := config.CredentialsJSON
	var app *firebase.App
	var err error

	if len(credentialsJSON) > 0 {
		// Initialize with service account credentials JSON
		opt := option.WithCredentialsJSON(credentialsJSON)
		conf := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, conf, opt)
	} else {
		// Use application default credentials (for local dev or GCP environment)
		conf := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, conf)
	}

	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting firestore client: %w", err)
	}

	return client, nil
}
