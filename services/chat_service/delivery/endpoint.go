package delivery

import (
	deliveryDao "go_ws/services/chat_service/delivery/dao/neo4j"
	"go_ws/services/chat_service/delivery/models"
	"go_ws/services/chat_service/delivery/usecases"
	"go_ws/shared/http_router"

	"net/http"
	"strings"
)

func DeliveriesByMessageHandler(findDeliveries deliveryDao.FindDeliveriesByMessageIdIn) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[map[string][]models.DeliveryView]

		execute(writer, func() (*map[string][]models.DeliveryView, error) {
			vars := http_router.Vars(request)
			chatIds := strings.Split(vars.Get("messageIds"), ",")

			result, err := findDeliveries(request.Context(), chatIds)

			return &result, err
		})
	}
}

func ConfirmMessageDeliveryStatusHandler(confirmDelivery usecases.ConfirmMessageDelivery) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, usecases.CreateDeliveryParams]

		execute(writer, request, func(body usecases.CreateDeliveryParams) (*any, error) {
			return nil, confirmDelivery(request.Context(), body)
		})

		writer.WriteHeader(http.StatusCreated)
	}
}
