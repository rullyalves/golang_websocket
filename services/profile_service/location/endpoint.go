package location

import (
	locationDao "go_ws/services/profile_service/location/dao/neo4j"
	"go_ws/services/profile_service/location/models"
	"go_ws/services/profile_service/location/usecases"

	"go_ws/shared/http_router"
	"net/http"
)

func FindLocationByProfileIdHandler(findLocationByProfileId locationDao.FindLocationByProfileId) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[models.LocationView]

		execute(writer, func() (*models.LocationView, error) {
			params := http_router.Params(request)
			id := params["userId"].(string)

			location, err := findLocationByProfileId(request.Context(), id)
			return location, err
		})
	}
}

func SaveLocationHandler(createLocation usecases.UpdateProfileLocation) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, usecases.CreateLocationByLocationParams]

		execute(writer, request, func(body usecases.CreateLocationByLocationParams) (*any, error) {

			params := http_router.Params(request)
			userId, _ := params["userId"].(string)
			return nil, createLocation(request.Context(), userId, body)
		})

		writer.WriteHeader(http.StatusNoContent)
	}
}
