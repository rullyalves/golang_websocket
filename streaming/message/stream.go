package message

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	messageDao "go_ws/streaming/message/dao"
	"go_ws/streaming/message/models"
	"go_ws/streaming/shared/broadcast"
)

func GetMessageStream(mongoDriver *mongo.Client) broadcast.Observer[*models.MessageEvent] {
	ctx := context.Background()

	stream := messageDao.FindAsStream(mongoDriver)

	dataStream, _ := stream(ctx)

	return dataStream
}
