package dislike

import (
	"github.com/google/uuid"
	dislikeDao "go_ws/services/relationship_service/dislike/dao/neo4j"
	"go_ws/services/relationship_service/dislike/usecases"
	"go_ws/shared/http_router"
	"net/http"
	"time"
)

type CreateDislikeParams struct {
	ReceiverID string `json:"receiverId" validate:"required,uuid4"`
}

func SendDislikeHandler(sendDislike usecases.SendDislike) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, CreateDislikeParams]

		execute(writer, request, func(body CreateDislikeParams) (*any, error) {
			params := http_router.Params(request)
			userId := params["userId"].(string)
			return nil, sendDislike(request.Context(), userId, dislikeDao.CreateDislikeDataParams{
				ID:         uuid.NewString(),
				ReceiverID: body.ReceiverID,
				CreatedAt:  time.Now().UTC(),
			})
		})

		writer.WriteHeader(http.StatusCreated)
	}
}
