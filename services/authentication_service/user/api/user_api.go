package api

import (
	"context"
	"firebase.google.com/go/v4/auth"
)

type CreateUserDataParams struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type CreateFirebaseUser func(context.Context, CreateUserDataParams) error

func CreateUser(authClient *auth.Client) CreateFirebaseUser {
	return func(ctx context.Context, params CreateUserDataParams) error {
		newUser := (&auth.UserToCreate{}).
			UID(params.ID).
			PhoneNumber(params.Username).
			Disabled(false)

		_, err := authClient.CreateUser(ctx, newUser)

		return err
	}
}
