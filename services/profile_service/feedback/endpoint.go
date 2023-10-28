package feedback

import (
	"github.com/google/uuid"
	feedbackDao "go_ws/services/profile_service/feedback/dao/neo4j"
	"go_ws/shared/http_router"
	"net/http"
	"time"
)

type CreateFeedbackParams struct {
	Description string `json:"description" validate:"required, min=10"`
}

func SendFeedbackHandler(createFeedback feedbackDao.SaveFeedback) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, CreateFeedbackParams]

		execute(writer, request, func(body CreateFeedbackParams) (*any, error) {
			params := http_router.Params(request)
			userId, _ := params["userId"].(string)
			return nil, createFeedback(request.Context(), userId, feedbackDao.CreateFeedbackDataParams{
				ID:          uuid.NewString(),
				Description: body.Description,
				CreatedAt:   time.Now().UTC(),
			})
		})

		writer.WriteHeader(http.StatusCreated)
	}
}
