package profile

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	profileDao "go_ws/streaming/profile/dao"
	"go_ws/streaming/profile/models"
	"go_ws/streaming/shared/broadcast"
)

func GetProfileStream(mongoDriver *mongo.Client) broadcast.Observer[*models.ProfileEvent] {
	ctx := context.Background()

	stream := profileDao.FindAsStream(mongoDriver)

	dataStream, _ := stream(ctx)

	return dataStream
}
