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
// Note: The base path (e.g., "images/") is handled by the gateway layer
func generateObjectKey(slug string) string {
	// Format: {slug}-{timestamp}.jpg (base path is added by gateway)
	return fmt.Sprintf("%s-%d.jpg", normalizeSlug(slug), time.Now().Unix())
}

// extractKeyFromURL extracts the object storage key from a full URL
func extractKeyFromURL(url string) string {
	// Extract everything after the bucket name
	// Example: https://storage.googleapis.com/bucket/images/sunset-123.jpg -> images/sunset-123.jpg
	// Example: https://storage.googleapis.com/bucket/sunset-123.jpg -> sunset-123.jpg
	parts := strings.Split(url, "/")
	if len(parts) >= 5 {
		// URL format: https://storage.googleapis.com/{bucket}/{path...}
		// We want everything from index 4 onwards
		return strings.Join(parts[4:], "/")
	}
	// Fallback: just get the filename
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
