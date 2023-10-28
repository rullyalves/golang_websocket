package authentication_service

import (
	"firebase.google.com/go/v4/auth"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/authentication_service/user"
	"go_ws/services/authentication_service/user/usecases"
	"go_ws/shared/http_router"
)

func Handle(router http_router.Router, neo4jDriver *neo4j.DriverWithContext, authClient *auth.Client) {
	findCurrentUser := usecases.NewFindCurrentUser(neo4jDriver)
	saveUser := usecases.NewCreateUser(authClient, neo4jDriver)
	router.Get("/users/current", user.FindCurrentUserHandler(findCurrentUser))
	router.Post("/users", user.SaveUserHandler(saveUser))
}
