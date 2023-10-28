package dislike

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	dislikeDao "go_ws/streaming/dislike/dao"
	"go_ws/streaming/dislike/models"
	"go_ws/streaming/shared/broadcast"
)

func GetDislikeStream(mongoDriver *mongo.Client) broadcast.Observer[*models.DislikeEvent] {
	ctx := context.Background()

	stream := dislikeDao.FindAsStream(mongoDriver)

	dataStream, _ := stream(ctx)

	return dataStream
}
