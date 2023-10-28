package block

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	blockDao "go_ws/streaming/block/dao"
	"go_ws/streaming/block/models"
	"go_ws/streaming/shared/broadcast"
)

func GetBlockStream(mongoDriver *mongo.Client) broadcast.Observer[*models.BlockEvent] {
	ctx := context.Background()

	stream := blockDao.FindAsStream(mongoDriver)

	dataStream, _ := stream(ctx)

	return dataStream
}
