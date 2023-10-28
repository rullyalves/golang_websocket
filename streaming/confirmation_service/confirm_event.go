package confirmation_service

import (
	"go_ws/shared/http_router"
	eventDao "go_ws/streaming/confirmation_service/dao/mongodb"
	"go_ws/streaming/confirmation_service/usecases"

	"net/http"
)

func ConfirmEventHandler(confirmEvent usecases.ConfirmSubscriptionEvent) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteVoidMutationRequest[any]

		execute(writer, func() (*any, error) {
			params := http_router.Params(request)
			vars := http_router.Vars(request)

			eventId := params["eventId"].(string)
			event := eventDao.Event(vars.Get("type"))

			return nil, confirmEvent(request.Context(), eventId, event)
		})

		writer.WriteHeader(http.StatusNoContent)
	}
}
