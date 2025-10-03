package middleware

import (
	"log"
	"net/http"
)

// Recovery middleware recovers from panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(r.Context())
				log.Printf("PANIC: %v request_id=%s", err, requestID)

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
