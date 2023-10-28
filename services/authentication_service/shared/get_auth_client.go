package shared

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func GetAuthClient(ctx context.Context, app *firebase.App) (*auth.Client, error) {

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
