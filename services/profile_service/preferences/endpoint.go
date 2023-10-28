package preferences

import (
	preferencesDao "go_ws/services/profile_service/preferences/dao/neo4j"
	"go_ws/services/profile_service/preferences/models"
	"go_ws/services/profile_service/preferences/usecases"
	"go_ws/shared/http_router"
	"net/http"
)

func FindPreferencesByProfileIdInHandler(findPreferencesByProfileId preferencesDao.FindPreferencesByProfileId) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[models.PreferencesView]

		execute(writer, func() (*models.PreferencesView, error) {
			params := http_router.Params(request)
			id := params["userId"].(string)

			preferences, err := findPreferencesByProfileId(request.Context(), id)
			return preferences, err
		})
	}
}

func UpdatePreferencesHandler(updatePreferences usecases.UpdateProfilePreferences) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, usecases.UpdatePreferencesData]

		execute(writer, request, func(body usecases.UpdatePreferencesData) (*any, error) {
			params := http_router.Params(request)
			userId := params["preferencesId"].(string)
			return nil, updatePreferences(request.Context(), userId, body)
		})

		writer.WriteHeader(http.StatusCreated)
	}
}
