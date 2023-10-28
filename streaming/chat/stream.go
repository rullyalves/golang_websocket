package chat

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	chatDao "go_ws/streaming/chat/dao"
	"go_ws/streaming/chat/models"
	"go_ws/streaming/shared/broadcast"
)

func GetChatStream(mongoDriver *mongo.Client) broadcast.Observer[*models.ChatEvent] {
	ctx := context.Background()

	stream := chatDao.FindAsStream(mongoDriver)

	dataStream, _ := stream(ctx)

	return dataStream
}
