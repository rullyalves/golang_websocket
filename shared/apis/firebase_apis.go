package apis

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"go_ws/services/authentication_service/shared"
	pushNotifications "go_ws/streaming/push_notifications"
	"google.golang.org/api/option"
	"log"
	"os"
)

func StartFirebaseApis(ctx context.Context) (*messaging.Client, *auth.Client, error) {
	googleServices := os.Getenv("FIREBASE_CONFIG")

	opt := option.WithCredentialsFile(googleServices)

	app, err := firebase.NewApp(ctx, nil, opt)

	if err != nil {
		log.Printf("error while starting %v\n", err)
		return nil, nil, err
	}

	messagingClient, err := pushNotifications.GetMessagingClient(ctx, app)

	if err != nil {
		log.Printf("error while starting %v\n", err)
		return nil, nil, err
	}

	authClient, err := shared.GetAuthClient(ctx, app)

	if err != nil {
		log.Printf("error while starting %v\n", err)
		return nil, nil, err
	}

	return messagingClient, authClient, nil
}
