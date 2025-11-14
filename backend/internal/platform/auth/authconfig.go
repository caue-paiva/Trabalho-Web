package auth

import (
	"firebase.google.com/go/v4/auth"
)

type AuthLevel int

const (
	AuthRequired AuthLevel = iota
	AuthOptional
)

type AuthConfig struct {
	Client *auth.Client
	Level  AuthLevel
}
