package messaging

import (
	"context"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
)

func GetMessagingClient(ctx context.Context, app *firebase.App) (*messaging.Client, error) {

	client, err := app.Messaging(ctx)

	if err != nil {
		return nil, err
	}

	return client, nil
}
