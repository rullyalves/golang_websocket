package filter

import (
	filterDao "go_ws/services/profile_service/filter/dao/neo4j"
	"go_ws/services/profile_service/filter/models"
	"go_ws/shared/http_router"
	"net/http"
)

func FindFilterByPreferencesIdInHandler(findFilterByProfileId filterDao.FindFilterByPreferencesId) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[[]models.FilterView]

		execute(writer, func() (*[]models.FilterView, error) {
			params := http_router.Params(request)
			id := params["preferencesId"].(string)

			filters, err := findFilterByProfileId(request.Context(), id)
			return filters, err
		})
	}
}
