package like

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	likeDao "go_ws/streaming/like/dao"
	"go_ws/streaming/like/models"
	"go_ws/streaming/shared/broadcast"
)

func GetLikeStream(mongoDriver *mongo.Client) broadcast.Observer[*models.LikeEvent] {
	ctx := context.Background()

	stream := likeDao.FindAsStream(mongoDriver)

	dataStream, _ := stream(ctx)

	return dataStream
}
