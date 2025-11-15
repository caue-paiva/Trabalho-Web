package auth

import (
	"strings"

	"firebase.google.com/go/v4/auth"
)

type AuthLevel int

const (
	AuthRequired AuthLevel = iota
	AuthOptional
)

func (l AuthLevel) String() string {
	switch l {

	case AuthRequired:
		return "required"
	case AuthOptional:
		return "optional"
	default:
		return "unknown"
	}
}

func AuthLevelFromString(value string) AuthLevel {
	lower := strings.ToLower(value)
	switch lower {

	case "required":
		return AuthRequired
	case "optional":
		return AuthOptional
	default:
		return AuthOptional
	}
}

// AuthConfig determines the configuration used for endpoint authorization in the server
type AuthConfig struct {
	Client *auth.Client
	Level  AuthLevel
}
