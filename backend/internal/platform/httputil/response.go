package httputil

import (
	"encoding/json"
	"net/http"

	customerrors "backend/internal/platform/errors"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// JSON writes a JSON response
func JSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// Log error but don't fail response
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

// Error writes an error response
func Error(w http.ResponseWriter, err error, status int) {
	JSON(w, ErrorResponse{Error: err.Error()}, status)
}

// ErrorFromDomain writes an error response using domain error mapping
func ErrorFromDomain(w http.ResponseWriter, err error) {
	status := customerrors.HTTPStatusFromError(err)
	Error(w, err, status)
}

// NoContent writes a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
