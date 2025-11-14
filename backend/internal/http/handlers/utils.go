package handlers

import "net/http"

type HandlerOption func()

// extractPathParam extracts a path parameter from the URL
// Uses Go 1.22+ PathValue method
func extractPathParam(r *http.Request, param string) string {
	return r.PathValue(param)
}
