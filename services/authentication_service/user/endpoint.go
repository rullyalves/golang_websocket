package user

import (
	"go_ws/services/authentication_service/user/models"
	"go_ws/services/authentication_service/user/usecases"
	"go_ws/shared/http_router"
	"net/http"
)

func FindCurrentUserHandler(findCurrentUser usecases.FindCurrentUser) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[models.UserView]

		execute(writer, func() (*models.UserView, error) {
			userId := request.Context().Value("userId").(string)
			return findCurrentUser(request.Context(), userId)
		})
	}
}

func SaveUserHandler(saveUser usecases.CreateUser) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, usecases.CreateUserParams]

		execute(writer, request, func(body usecases.CreateUserParams) (*any, error) {
			return nil, saveUser(request.Context(), body)
		})

		writer.WriteHeader(http.StatusNoContent)
	}
}
