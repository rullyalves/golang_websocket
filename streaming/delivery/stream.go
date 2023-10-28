package delivery

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	deliveryDao "go_ws/streaming/delivery/dao"
	"go_ws/streaming/delivery/models"
	"go_ws/streaming/shared/broadcast"
)

func GetDeliveryStream(mongoDriver *mongo.Client) broadcast.Observer[*models.DeliveryEvent] {
	ctx := context.Background()

	stream := deliveryDao.FindAsStream(mongoDriver)

	dataStream, _ := stream(ctx)

	return dataStream
}
