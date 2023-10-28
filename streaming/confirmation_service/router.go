package confirmation_service

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go_ws/shared/http_router"
	eventDao "go_ws/streaming/confirmation_service/dao/mongodb"
	"go_ws/streaming/confirmation_service/usecases"
)

func Handle(router http_router.Router, mongoDriver *mongo.Client) {
	deleteEvent := eventDao.DeleteById(mongoDriver)
	confirmEvent := usecases.ConfirmEvent(deleteEvent)
	router.Delete("/events/:eventId", ConfirmEventHandler(confirmEvent))
}
