package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger middleware logs HTTP requests
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Log request details
		duration := time.Since(start)
		requestID := GetRequestID(r.Context())

		log.Printf("[%s] %s %s %d %v request_id=%s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			wrapped.status,
			duration,
			requestID,
		)
	})
}
