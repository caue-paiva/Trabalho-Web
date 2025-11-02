package server

import (
	"fmt"
	"strings"
	"time"
)

// normalizeSlug normalizes a slug by lowercasing, trimming, and replacing spaces with hyphens
func normalizeSlug(slug string) string {
	normalized := strings.TrimSpace(strings.ToLower(slug))
	normalized = strings.ReplaceAll(normalized, " ", "-")
	return normalized
}

// generateObjectKey generates a unique object storage key for an image
func generateObjectKey(slug string) string {
	// Format: images/{slug}-{timestamp}.jpg
	return fmt.Sprintf("images/%s-%d.jpg", normalizeSlug(slug), time.Now().Unix())
}

// extractKeyFromURL extracts the object storage key from a full URL
func extractKeyFromURL(url string) string {
	// Simple extraction: assume URL ends with the key
	// Example: https://storage.googleapis.com/bucket/images/sunset-123.jpg -> images/sunset-123.jpg
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return ""
}
