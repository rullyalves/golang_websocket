package usecases

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/authentication_service/user/api"
	userDao "go_ws/services/authentication_service/user/dao"
	"time"
)

type CreateUserParams struct {
	Username string `json:"username" validate:"required,e164"`
}

type CreateUser func(ctx context.Context, params CreateUserParams) error

func NewCreateUser(authClient *auth.Client, neo4jDriver *neo4j.DriverWithContext) CreateUser {
	saveFirebaseUser := api.CreateUser(authClient)
	saveUser := userDao.Save(neo4jDriver)
	findByUsername := userDao.FindByUsername(neo4jDriver)
	return createUser(saveFirebaseUser, findByUsername, saveUser)
}

func createUser(saveFirebaseUser api.CreateFirebaseUser, findByUsername userDao.FindUserByUsername, saveUser userDao.SaveUser) CreateUser {
	return func(ctx context.Context, params CreateUserParams) error {

		username := params.Username

		result, err := findByUsername(ctx, username)

		if err != nil {
			return err
		}

		if result != nil {
			return nil
		}

		userId := uuid.NewString()

		err = saveFirebaseUser(ctx, api.CreateUserDataParams{
			ID:       userId,
			Username: username,
		})

		if err != nil && !auth.IsPhoneNumberAlreadyExists(err) {
			return err
		}

		if err != nil {
			return err
		}

		_, err = saveUser(ctx, userDao.CreateUserDataParams{
			ID:        userId,
			Username:  username,
			CreatedAt: time.Now().UTC(),
		})

		return err
	}
}
