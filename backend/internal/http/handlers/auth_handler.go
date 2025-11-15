package handlers

import (
	"net/http"

	"backend/internal/platform/httputil"
)

// Authorized handles GET /authorized
// This endpoint is protected by NewForceAuthMiddlewareFunc and only returns
// a success response if the request is properly authenticated.
func (h *BaseHandler) Authorized(w http.ResponseWriter, r *http.Request) {
	// If we reach this handler, the authentication middleware has already
	// verified that the request is authorized (otherwise it would have
	// returned 401 Unauthorized)
	response := map[string]interface{}{
		"authorized": true,
		"message":    "Request is authorized",
		"status":     "success",
	}

	httputil.JSON(w, response, http.StatusOK)
}
