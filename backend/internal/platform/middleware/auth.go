package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	authcfg "backend/internal/platform/auth"
)

// Middleware is a func which takes in an http request and returns it and an error
type Middleware func(*http.Request) (*http.Request, error)

func NewAuthMiddlewareFunc(nextHandle func(w http.ResponseWriter, r *http.Request), authCfg authcfg.AuthConfig, logger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	if authCfg.Client == nil {
		return func(w http.ResponseWriter, r *http.Request) {
			nextHandle(w, r)
		}
	}

	switch authCfg.Level {

	case authcfg.AuthRequired:
		return func(w http.ResponseWriter, r *http.Request) {
			if r == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// log token extraction and enforce it
			idToken, err := getIdToken(r)
			if err != nil {
				logTokenNotFound(logger, r, err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			logTokenFound(logger, r, idToken)

			// Verify the token
			_, err = authCfg.Client.VerifyIDToken(r.Context(), idToken)
			if err != nil {
				logTokenVerificationFailed(logger, r, idToken, err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			nextHandle(w, r)
		}
	default:
		return func(w http.ResponseWriter, r *http.Request) {
			// log token extraction but do not enforce it
			idToken, err := getIdToken(r)
			if err != nil {
				logTokenNotFound(logger, r, err)
			} else {
				logTokenFound(logger, r, idToken)
			}
			nextHandle(w, r)
		}
	}

}

// getRequestOrigin gets the origin from the request (Origin header or RemoteAddr)
func getRequestOrigin(r *http.Request) string {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = r.RemoteAddr
	}
	return origin
}

// logTokenNotFound logs when a token is not found
func logTokenNotFound(logger *log.Logger, r *http.Request, err error) {
	origin := getRequestOrigin(r)
	authHeader := r.Header.Get("Authorization")
	logger.Printf("[AUTH] Token not found - origin=%s auth_header=%s error=%v",
		origin,
		authHeader,
		err,
	)
}

// logTokenFound logs when a token is found (first 10 chars only for security)
func logTokenFound(logger *log.Logger, r *http.Request, idToken string) {
	origin := getRequestOrigin(r)
	authHeader := r.Header.Get("Authorization")
	tokenPrefix := idToken
	if len(tokenPrefix) > 10 {
		tokenPrefix = tokenPrefix[:10]
	}
	logger.Printf("[AUTH] Token found - origin=%s auth_header=%s token_prefix=%s",
		origin,
		authHeader,
		tokenPrefix,
	)
}

// logTokenVerificationFailed logs when token verification fails
func logTokenVerificationFailed(logger *log.Logger, r *http.Request, idToken string, err error) {
	origin := getRequestOrigin(r)
	tokenPrefix := idToken
	if len(tokenPrefix) > 10 {
		tokenPrefix = tokenPrefix[:10]
	}
	logger.Printf("[AUTH] Token verification failed - origin=%s token_prefix=%s error=%v",
		origin,
		tokenPrefix,
		err,
	)
}

func getIdToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("malformed authorization header")
	}

	idToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if idToken == "" {
		return "", errors.New("id token not found")
	}

	return idToken, nil
}
