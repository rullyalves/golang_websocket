package profile_service

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/profile_service/feedback"
	feedbackDao "go_ws/services/profile_service/feedback/dao/neo4j"
	"go_ws/services/profile_service/filter"
	filterDao "go_ws/services/profile_service/filter/dao/neo4j"
	"go_ws/services/profile_service/location"
	locationDao "go_ws/services/profile_service/location/dao/neo4j"
	locationUseCases "go_ws/services/profile_service/location/usecases"
	"go_ws/services/profile_service/options"
	optionDao "go_ws/services/profile_service/options/dao/neo4j"
	"go_ws/services/profile_service/preferences"
	preferencesDao "go_ws/services/profile_service/preferences/dao/neo4j"
	preferencesUseCases "go_ws/services/profile_service/preferences/usecases"
	"go_ws/services/profile_service/profile"
	profileDao "go_ws/services/profile_service/profile/dao/neo4j"
	profileUseCases "go_ws/services/profile_service/profile/usecases"
	"go_ws/services/profile_service/settings"
	settingsDao "go_ws/services/profile_service/settings/dao/neo4j"
	"go_ws/shared/http_router"
)

func Handle(router http_router.Router, neo4jDriver *neo4j.DriverWithContext, sqsClient *sqs.Client) {

	// profiles
	findByIdIn := profileDao.FindByIdIn(neo4jDriver)
	findInterestProfiles := profileDao.FindInterestProfilesIn(neo4jDriver)
	saveProfile := profileUseCases.NewCreateProfile(neo4jDriver)
	updateProfile := profileUseCases.NewUpdateProfile(neo4jDriver, sqsClient)

	router.Get("/profiles", profile.FindProfilesByIdInHandler(findByIdIn))
	router.Get("/profiles/:userId/interest", profile.FindInterestProfilesHandler(findInterestProfiles))
	router.Post("/profiles", profile.SaveProfileHandler(saveProfile))
	router.Patch("/profiles/:userId", profile.UpdateProfileHandler(updateProfile))

	// location
	findLocation := locationDao.FindProfileId(neo4jDriver)
	updateLocation := locationUseCases.NewUpdateLocation(neo4jDriver)
	router.Get("/profiles/:userId/location", location.FindLocationByProfileIdHandler(findLocation))
	router.Post("/profiles/:userId/location", location.SaveLocationHandler(updateLocation))

	// preferences
	findByProfile := preferencesDao.FindByProfileId(neo4jDriver)
	updatePreferences := preferencesUseCases.NewUpdatePreferences(neo4jDriver)

	router.Get("/profiles/:userId/preferences", preferences.FindPreferencesByProfileIdInHandler(findByProfile))
	router.Patch("/preferences/:preferencesId", preferences.UpdatePreferencesHandler(updatePreferences))

	// filters
	findFilterByPreferences := filterDao.FindByPreferencesId(neo4jDriver)
	router.Get("/preferences/:preferencesId/filters", filter.FindFilterByPreferencesIdInHandler(findFilterByPreferences))

	// options
	findOptions := optionDao.FindAll(neo4jDriver)
	saveOption := optionDao.Save(neo4jDriver)
	deleteOption := optionDao.Delete(neo4jDriver)

	router.Get("/options", options.FindOptionsHandler(findOptions))
	router.Post("/options", options.SaveOptionHandler(saveOption))
	router.Delete("/options", options.DeleteOptionHandler(deleteOption))

	// feedback
	sendFeedback := feedbackDao.Save(neo4jDriver)
	router.Post("/users/:userId/feedback", feedback.SendFeedbackHandler(sendFeedback))

	// settings
	findSettings := settingsDao.FindProfileIdIn(neo4jDriver)
	saveSettings := settingsDao.Save(neo4jDriver)
	router.Get("/users/:userId/settings", settings.FindSettingsHandler(findSettings))
	router.Post("/users/:userId/settings", settings.SaveSettingsHandler(saveSettings))
}
