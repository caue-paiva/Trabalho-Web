package clients

import (
	"context"
	"errors"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func NewFirebaseAuthClient(ctx context.Context, fbApp *firebase.App) (*auth.Client, error) {
	if fbApp == nil {
		return nil, errors.New("firebase app client cant be nil")
	}

	return fbApp.Auth(ctx)
}
