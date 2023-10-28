package usecases

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	userDao "go_ws/services/authentication_service/user/dao"
	"go_ws/services/authentication_service/user/models"
)

type FindCurrentUser func(context.Context, string) (*models.UserView, error)

func NewFindCurrentUser(neo4jDriver *neo4j.DriverWithContext) FindCurrentUser {
	findUser := userDao.FindByUserIdIn(neo4jDriver)
	return FindUser(findUser)
}

func FindUser(findByIdIn userDao.FindUserByIdIn) FindCurrentUser {

	return func(ctx context.Context, userId string) (*models.UserView, error) {

		result, err := findByIdIn(ctx, []string{userId})

		if len(result) == 0 {
			return nil, err
		}

		return &result[0], nil
	}
}
