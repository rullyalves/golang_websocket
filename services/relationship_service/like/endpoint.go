package like

import (
	likeDao "go_ws/services/relationship_service/like/dao/neo4j"
	"go_ws/services/relationship_service/like/models"
	"go_ws/services/relationship_service/like/usecases"
	"go_ws/shared/http_router"
	"net/http"
)

func FindReceivedLikesHandler(findReceivedLikes likeDao.FindReceivedLikesByProfileId) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[[]models.LikeView]

		execute(writer, func() (*[]models.LikeView, error) {
			params := http_router.Params(request)
			userId := params["userId"].(string)

			like, err := findReceivedLikes(request.Context(), userId)

			if err != nil {
				return nil, err
			}

			return like, err
		})

		writer.WriteHeader(http.StatusOK)
	}
}

func SendLikeHandler(sendLike usecases.SendLike) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, usecases.CreateLikeParams]

		execute(writer, request, func(body usecases.CreateLikeParams) (*any, error) {
			params := http_router.Params(request)
			userId := params["userId"].(string)
			return nil, sendLike(request.Context(), userId, body)
		})

		writer.WriteHeader(http.StatusCreated)
	}
}
