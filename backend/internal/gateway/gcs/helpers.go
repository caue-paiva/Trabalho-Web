package gcs

import (
	"mime"
	"path/filepath"
)

// detectContentType returns the appropriate Content-Type for a file based on its extension
func detectContentType(key string) string {
	ext := filepath.Ext(key)

	// Try standard mime type detection first
	contentType := mime.TypeByExtension(ext)
	if contentType != "" {
		return contentType
	}

	// Fallback to explicit mapping for common image types
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".bmp":
		return "image/bmp"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}
