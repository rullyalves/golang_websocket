package options

import (
	"github.com/google/uuid"
	optionDao "go_ws/services/profile_service/options/dao/neo4j"
	"go_ws/services/profile_service/options/models"
	"go_ws/shared/http_router"
	"net/http"
	"time"
)

type CreateOptionParams struct {
	Name      string            `json:"name" validate:"required"`
	Type      models.OptionType `json:"type" validate:"required,lowercase"`
	CreatedAt time.Time         `json:"createdAt" validate:"required,datetime"`
}

func FindOptionsHandler(findOptions optionDao.FindOptions) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[[]models.OptionView]

		execute(writer, func() (*[]models.OptionView, error) {
			result, err := findOptions(request.Context())
			return &result, err
		})

		writer.WriteHeader(http.StatusCreated)
	}
}

func SaveOptionHandler(saveOption optionDao.SaveOption) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, CreateOptionParams]

		execute(writer, request, func(body CreateOptionParams) (*any, error) {
			err := saveOption(request.Context(), optionDao.CreateOptionDataParams{
				ID:        uuid.NewString(),
				Name:      body.Name,
				Type:      body.Type,
				CreatedAt: time.Now().UTC(),
			})
			return nil, err
		})

		writer.WriteHeader(http.StatusCreated)
	}
}

func DeleteOptionHandler(deleteOption optionDao.DeleteOptionById) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteVoidMutationRequest[any]

		execute(writer, func() (*any, error) {
			params := http_router.Params(request)
			id := params["id"].(string)

			return nil, deleteOption(request.Context(), id)
		})

		writer.WriteHeader(http.StatusCreated)
	}
}
