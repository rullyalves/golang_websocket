package blocked_user

import (
	"github.com/google/uuid"
	blockEventDao "go_ws/services/relationship_service/blocked_user/dao/mongo"
	blockDao "go_ws/services/relationship_service/blocked_user/dao/neo4j"
	"go_ws/shared/http_router"
	"net/http"
	"time"
)

type CreateBlockParams struct {
	ReceiverID  string  `json:"receiverId" validate:"required,uuid4"`
	ComplaintID *string `json:"complaintId" validate:"omitempty,uuid4"`
}

func SendBlockHandler(saveBlock blockDao.SaveBlock) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, CreateBlockParams]

		execute(writer, request, func(body CreateBlockParams) (*any, error) {
			params := http_router.Params(request)
			userId, _ := params["userId"].(string)
			return nil, saveBlock(request.Context(), userId, blockDao.CreateBlockDataParams{
				ID:          uuid.NewString(),
				ReceiverID:  body.ReceiverID,
				ComplaintID: body.ComplaintID,
				CreatedAt:   time.Now().UTC(),
			})
		})

		writer.WriteHeader(http.StatusCreated)
	}
}

func ConfirmBlockEventHandler(deleteEvent blockEventDao.DeleteBlockEventById) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteVoidMutationRequest[any]

		execute(writer, func() (*any, error) {
			vars := request.URL.Query()
			id := vars.Get("eventId")
			return nil, deleteEvent(request.Context(), id)
		})

		writer.WriteHeader(http.StatusNoContent)
	}
}
