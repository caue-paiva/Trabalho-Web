package clients

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"backend/configs"
)

func NewFirebaseAppClient(ctx context.Context, cfg configs.FirebaseConfig) (*firebase.App, error) {
	fbCfg := &firebase.Config{}
	if cfg.ProjectID != "" {
		fbCfg.ProjectID = cfg.ProjectID
	}

	var app *firebase.App
	var err error
	if len(cfg.CredentialsJSON) > 0 {
		app, err = firebase.NewApp(ctx, fbCfg, option.WithCredentialsJSON(cfg.CredentialsJSON))
	} else {
		app, err = firebase.NewApp(ctx, fbCfg)
	}

	return app, err
}
