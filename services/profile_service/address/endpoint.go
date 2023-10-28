package address

import (
	"go_ws/services/profile_service/address/usecases"
	"go_ws/shared/http_router"
	"net/http"
)

func SaveAddressHandler(createAddress usecases.UpdateProfileLocation) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, usecases.CreateAddressByLocationParams]

		execute(writer, request, func(body usecases.CreateAddressByLocationParams) (*any, error) {

			params := http_router.Params(request)
			userId, _ := params["userId"].(string)
			return nil, createAddress(request.Context(), userId, body)
		})

		writer.WriteHeader(http.StatusNoContent)
	}
}
