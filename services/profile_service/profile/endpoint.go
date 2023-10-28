package profile

import (
	profileDao "go_ws/services/profile_service/profile/dao/neo4j"
	"go_ws/services/profile_service/profile/models"
	"go_ws/services/profile_service/profile/usecases"
	"go_ws/shared/http_router"
	"net/http"
	"strconv"
	"strings"
)

func FindProfilesByIdInHandler(findProfiles profileDao.FindProfileByIdIn) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[[]models.ProfileView]

		execute(writer, func() (*[]models.ProfileView, error) {
			vars := http_router.Vars(request)
			ids := strings.Split(vars.Get("ids"), ",")

			profiles, err := findProfiles(request.Context(), ids)

			return profiles, err
		})
	}
}

func FindInterestProfilesHandler(findInterestProfiles profileDao.FindInterestProfilesByProfileId) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[[]models.ProfileView]

		execute(writer, func() (*[]models.ProfileView, error) {
			params := http_router.Params(request)
			vars := http_router.Vars(request)
			userId := params["userId"].(string)
			limit, err := strconv.Atoi(vars.Get("limit"))

			if err != nil {
				return nil, err
			}

			profiles, err := findInterestProfiles(request.Context(), userId, limit)

			return profiles, err
		})
	}
}

func UpdateProfileHandler(createProfile usecases.UpdateProfileData) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, usecases.UpdateProfileParams]
		execute(writer, request, func(body usecases.UpdateProfileParams) (*any, error) {
			params := http_router.Params(request)
			userId := params["userId"].(string)
			return nil, createProfile(request.Context(), userId, body)
		})

		writer.WriteHeader(http.StatusNoContent)
	}
}

func SaveProfileHandler(createProfile usecases.CreateProfileData) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, usecases.CreateProfileParams]
		execute(writer, request, func(body usecases.CreateProfileParams) (*any, error) {
			return nil, createProfile(request.Context(), body)
		})

		writer.WriteHeader(http.StatusCreated)
	}
}
