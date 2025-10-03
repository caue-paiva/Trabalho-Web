package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// generateRequestID generates a random request ID using crypto/rand
func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to a simple timestamp-based ID if random generation fails
		return "fallback"
	}
	return hex.EncodeToString(b)
}

// RequestID middleware generates and injects a request ID into the context
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate request ID
		requestID := generateRequestID()

		// Inject into context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Add to response header
		w.Header().Set("X-Request-ID", requestID)

		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
