package settings

import (
	"github.com/google/uuid"
	settingsDao "go_ws/services/profile_service/settings/dao/neo4j"
	"go_ws/services/profile_service/settings/models"
	"go_ws/shared/http_router"
	"net/http"
)

type CreateSettingsParams struct {
	IsVisible                      *bool `json:"isVisible"`
	AllowReceiveMatchNotifications *bool `json:"allowReceiveMatchNotifications"`
	AllowReceiveLikeNotifications  *bool `json:"allowReceiveLikeNotifications"`
}

func FindSettingsHandler(findSettings settingsDao.FindSettingsByProfileIdIn) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[models.SettingsView]

		execute(writer, func() (*models.SettingsView, error) {
			params := http_router.Params(request)
			userId, _ := params["userId"].(string)

			result, err := findSettings(request.Context(), []string{userId})
			//TODO: melhorar depois

			return &result[0].Settings, err
		})

		writer.WriteHeader(http.StatusOK)
	}
}

func SaveSettingsHandler(createSettings settingsDao.SaveSettings) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, CreateSettingsParams]

		execute(writer, request, func(body CreateSettingsParams) (*any, error) {
			params := http_router.Params(request)
			userId, _ := params["userId"].(string)
			return nil, createSettings(request.Context(), userId, settingsDao.CreateSettingsDataParams{
				ID:                             uuid.NewString(),
				IsVisible:                      body.IsVisible,
				AllowReceiveLikeNotifications:  body.AllowReceiveMatchNotifications,
				AllowReceiveMatchNotifications: body.AllowReceiveMatchNotifications,
			})
		})

		writer.WriteHeader(http.StatusCreated)
	}
}
