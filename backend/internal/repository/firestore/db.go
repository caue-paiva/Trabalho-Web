package firestore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"backend/internal/entities"
	"backend/internal/server"
)

// Compile-time check that DBRepository implements server.DBPort
var _ server.DBPort = (*DBRepository)(nil)

// DBRepository implements server.DBPort using Firestore
type DBRepository struct {
	client      *firestore.Client
	collections CollectionNames
}

// NewDBRepository creates a new Firestore DB repository
func NewDBRepository(client *firestore.Client, collections CollectionNames) *DBRepository {
	return &DBRepository{
		client:      client,
		collections: collections,
	}
}

// =======================
// TEXT OPERATIONS
// =======================

func (r *DBRepository) GetTextBySlug(ctx context.Context, slug string) (entities.Text, error) {
	iter := r.client.Collection(r.collections.Texts).Where("slug", "==", slug).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return entities.Text{}, fmt.Errorf("text with slug %s not found", slug)
	}
	if err != nil {
		return entities.Text{}, fmt.Errorf("error fetching text: %w", err)
	}

	var text entities.Text
	if err := doc.DataTo(&text); err != nil {
		return entities.Text{}, fmt.Errorf("error parsing text: %w", err)
	}
	text.ID = doc.Ref.ID
	return text, nil
}

func (r *DBRepository) GetTextByID(ctx context.Context, id string) (entities.Text, error) {
	doc, err := r.client.Collection(r.collections.Texts).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return entities.Text{}, fmt.Errorf("text with id %s not found", id)
		}
		return entities.Text{}, fmt.Errorf("error fetching text: %w", err)
	}

	var text entities.Text
	if err := doc.DataTo(&text); err != nil {
		return entities.Text{}, fmt.Errorf("error parsing text: %w", err)
	}
	text.ID = doc.Ref.ID
	return text, nil
}

func (r *DBRepository) GetTextsByPageID(ctx context.Context, pageID string) ([]entities.Text, error) {
	iter := r.client.Collection(r.collections.Texts).Where("pageID", "==", pageID).Documents(ctx)
	return r.textsFromIterator(iter)
}

func (r *DBRepository) ListTextsByPageSlug(ctx context.Context, pageSlug string) ([]entities.Text, error) {
	iter := r.client.Collection(r.collections.Texts).Where("pageSlug", "==", pageSlug).Documents(ctx)
	return r.textsFromIterator(iter)
}

func (r *DBRepository) ListAllTexts(ctx context.Context) ([]entities.Text, error) {
	iter := r.client.Collection(r.collections.Texts).Documents(ctx)
	return r.textsFromIterator(iter)
}

func (r *DBRepository) CreateText(ctx context.Context, text entities.Text) (entities.Text, error) {
	// Generate new document reference
	docRef := r.client.Collection(r.collections.Texts).NewDoc()
	text.ID = docRef.ID

	// Set timestamps if not already set
	if text.CreatedAt.IsZero() {
		text.CreatedAt = time.Now()
	}
	if text.UpdatedAt.IsZero() {
		text.UpdatedAt = time.Now()
	}

	// Create document
	if _, err := docRef.Set(ctx, text); err != nil {
		return entities.Text{}, fmt.Errorf("error creating text: %w", err)
	}

	return text, nil
}

func (r *DBRepository) UpdateText(ctx context.Context, id string, patch entities.Text) (entities.Text, error) {
	docRef := r.client.Collection(r.collections.Texts).Doc(id)

	// Update timestamp
	patch.UpdatedAt = time.Now()

	// Build update map (only update provided fields)
	updates := []firestore.Update{
		{Path: "updatedAt", Value: patch.UpdatedAt},
	}
	if patch.Content != "" {
		updates = append(updates, firestore.Update{Path: "content", Value: patch.Content})
	}
	if patch.Slug != "" {
		updates = append(updates, firestore.Update{Path: "slug", Value: patch.Slug})
	}
	if patch.PageID != "" {
		updates = append(updates, firestore.Update{Path: "pageId", Value: patch.PageID})
	}
	if patch.PageSlug != "" {
		updates = append(updates, firestore.Update{Path: "pageSlug", Value: patch.PageSlug})
	}
	if patch.LastUpdatedBy != "" {
		updates = append(updates, firestore.Update{Path: "lastUpdatedBy", Value: patch.LastUpdatedBy})
	}

	if _, err := docRef.Update(ctx, updates); err != nil {
		if status.Code(err) == codes.NotFound {
			return entities.Text{}, fmt.Errorf("text with id %s not found", id)
		}
		return entities.Text{}, fmt.Errorf("error updating text: %w", err)
	}

	// Fetch and return updated document
	return r.GetTextByID(ctx, id)
}

func (r *DBRepository) DeleteText(ctx context.Context, id string) error {
	if _, err := r.client.Collection(r.collections.Texts).Doc(id).Delete(ctx); err != nil {
		return fmt.Errorf("error deleting text: %w", err)
	}
	return nil
}

// =======================
// IMAGE OPERATIONS
// =======================

func (r *DBRepository) GetImageByID(ctx context.Context, id string) (entities.Image, error) {
	doc, err := r.client.Collection(r.collections.Images).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return entities.Image{}, fmt.Errorf("image with id %s not found", id)
		}
		return entities.Image{}, fmt.Errorf("error fetching image: %w", err)
	}

	var image entities.Image
	if err := doc.DataTo(&image); err != nil {
		return entities.Image{}, fmt.Errorf("error parsing image: %w", err)
	}
	image.ID = doc.Ref.ID
	return image, nil
}

func (r *DBRepository) GetImagesByGallerySlug(ctx context.Context, slug string) ([]entities.Image, error) {
	iter := r.client.Collection(r.collections.Images).Where("slug", "==", slug).Documents(ctx)
	return r.imagesFromIterator(iter)
}

func (r *DBRepository) CreateImageMeta(ctx context.Context, img entities.Image) (entities.Image, error) {
	// Generate new document reference
	docRef := r.client.Collection(r.collections.Images).NewDoc()
	img.ID = docRef.ID

	// Set timestamps if not already set
	if img.CreatedAt.IsZero() {
		img.CreatedAt = time.Now()
	}
	if img.UpdatedAt.IsZero() {
		img.UpdatedAt = time.Now()
	}

	// Create document
	if _, err := docRef.Set(ctx, img); err != nil {
		return entities.Image{}, fmt.Errorf("error creating image: %w", err)
	}

	return img, nil
}

func (r *DBRepository) UpdateImageMeta(ctx context.Context, id string, patch entities.Image) (entities.Image, error) {
	docRef := r.client.Collection(r.collections.Images).Doc(id)

	// Update timestamp
	patch.UpdatedAt = time.Now()

	// Build update map
	updates := []firestore.Update{
		{Path: "updatedAt", Value: patch.UpdatedAt},
	}
	if patch.Name != "" {
		updates = append(updates, firestore.Update{Path: "name", Value: patch.Name})
	}
	if patch.Text != "" {
		updates = append(updates, firestore.Update{Path: "text", Value: patch.Text})
	}
	if patch.Slug != "" {
		updates = append(updates, firestore.Update{Path: "slug", Value: patch.Slug})
	}
	if patch.ObjectURL != "" {
		updates = append(updates, firestore.Update{Path: "objectUrl", Value: patch.ObjectURL})
	}
	if patch.Location != "" {
		updates = append(updates, firestore.Update{Path: "location", Value: patch.Location})
	}
	if !patch.Date.IsZero() {
		updates = append(updates, firestore.Update{Path: "date", Value: patch.Date})
	}
	if patch.LastUpdatedBy != "" {
		updates = append(updates, firestore.Update{Path: "lastUpdatedBy", Value: patch.LastUpdatedBy})
	}

	if _, err := docRef.Update(ctx, updates); err != nil {
		if status.Code(err) == codes.NotFound {
			return entities.Image{}, fmt.Errorf("image with id %s not found", id)
		}
		return entities.Image{}, fmt.Errorf("error updating image: %w", err)
	}

	// Fetch and return updated document
	return r.GetImageByID(ctx, id)
}

func (r *DBRepository) DeleteImageMeta(ctx context.Context, id string) error {
	if _, err := r.client.Collection(r.collections.Images).Doc(id).Delete(ctx); err != nil {
		return fmt.Errorf("error deleting image: %w", err)
	}
	return nil
}

// =======================
// TIMELINE OPERATIONS
// =======================

func (r *DBRepository) GetTimelineEntryByID(ctx context.Context, id string) (entities.TimelineEntry, error) {
	doc, err := r.client.Collection(r.collections.TimelineEntries).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return entities.TimelineEntry{}, fmt.Errorf("timeline entry with id %s not found", id)
		}
		return entities.TimelineEntry{}, fmt.Errorf("error fetching timeline entry: %w", err)
	}

	var entry entities.TimelineEntry
	if err := doc.DataTo(&entry); err != nil {
		return entities.TimelineEntry{}, fmt.Errorf("error parsing timeline entry: %w", err)
	}
	entry.ID = doc.Ref.ID
	return entry, nil
}

func (r *DBRepository) ListTimelineEntries(ctx context.Context) ([]entities.TimelineEntry, error) {
	iter := r.client.Collection(r.collections.TimelineEntries).OrderBy("date", firestore.Asc).Documents(ctx)
	return r.timelineEntriesFromIterator(iter)
}

func (r *DBRepository) CreateTimelineEntry(ctx context.Context, entry entities.TimelineEntry) (entities.TimelineEntry, error) {
	// Generate new document reference
	docRef := r.client.Collection(r.collections.TimelineEntries).NewDoc()
	entry.ID = docRef.ID

	// Set timestamps if not already set
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	if entry.UpdatedAt.IsZero() {
		entry.UpdatedAt = time.Now()
	}

	// Create document
	if _, err := docRef.Set(ctx, entry); err != nil {
		return entities.TimelineEntry{}, fmt.Errorf("error creating timeline entry: %w", err)
	}

	return entry, nil
}

func (r *DBRepository) UpdateTimelineEntry(ctx context.Context, id string, patch entities.TimelineEntry) (entities.TimelineEntry, error) {
	docRef := r.client.Collection(r.collections.TimelineEntries).Doc(id)

	// Update timestamp
	patch.UpdatedAt = time.Now()

	// Build update map
	updates := []firestore.Update{
		{Path: "updatedAt", Value: patch.UpdatedAt},
	}
	if patch.Name != "" {
		updates = append(updates, firestore.Update{Path: "name", Value: patch.Name})
	}
	if patch.Text != "" {
		updates = append(updates, firestore.Update{Path: "text", Value: patch.Text})
	}
	if patch.Location != "" {
		updates = append(updates, firestore.Update{Path: "location", Value: patch.Location})
	}
	if !patch.Date.IsZero() {
		updates = append(updates, firestore.Update{Path: "date", Value: patch.Date})
	}
	if patch.LastUpdatedBy != "" {
		updates = append(updates, firestore.Update{Path: "lastUpdatedBy", Value: patch.LastUpdatedBy})
	}

	if _, err := docRef.Update(ctx, updates); err != nil {
		if status.Code(err) == codes.NotFound {
			return entities.TimelineEntry{}, fmt.Errorf("timeline entry with id %s not found", id)
		}
		return entities.TimelineEntry{}, fmt.Errorf("error updating timeline entry: %w", err)
	}

	// Fetch and return updated document
	return r.GetTimelineEntryByID(ctx, id)
}

func (r *DBRepository) DeleteTimelineEntry(ctx context.Context, id string) error {
	if _, err := r.client.Collection(r.collections.TimelineEntries).Doc(id).Delete(ctx); err != nil {
		return fmt.Errorf("error deleting timeline entry: %w", err)
	}
	return nil
}

// =======================
// HELPER METHODS
// =======================

func (r *DBRepository) textsFromIterator(iter *firestore.DocumentIterator) ([]entities.Text, error) {
	var texts []entities.Text
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating texts: %w", err)
		}

		var text entities.Text
		if err := doc.DataTo(&text); err != nil {
			continue // Skip malformed documents
		}
		text.ID = doc.Ref.ID
		texts = append(texts, text)
	}
	return texts, nil
}

func (r *DBRepository) imagesFromIterator(iter *firestore.DocumentIterator) ([]entities.Image, error) {
	var images []entities.Image
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating images: %w", err)
		}

		var image entities.Image
		if err := doc.DataTo(&image); err != nil {
			continue // Skip malformed documents
		}
		image.ID = doc.Ref.ID
		images = append(images, image)
	}
	return images, nil
}

func (r *DBRepository) timelineEntriesFromIterator(iter *firestore.DocumentIterator) ([]entities.TimelineEntry, error) {
	var entries []entities.TimelineEntry
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating timeline entries: %w", err)
		}

		var entry entities.TimelineEntry
		if err := doc.DataTo(&entry); err != nil {
			continue // Skip malformed documents
		}
		entry.ID = doc.Ref.ID
		entries = append(entries, entry)
	}
	return entries, nil
}
