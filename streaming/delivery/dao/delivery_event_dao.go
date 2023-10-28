package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go_ws/streaming/shared/broadcast"

	mongodb "go_ws/shared/database/mongo"
	"go_ws/streaming/delivery/models"
)

const (
	database   = "events"
	collection = "deliveries"
)

type FindDeliveriesAsStream func(ctx context.Context) (broadcast.Observer[*models.DeliveryEvent], error)

type FindDeliveriesByUserId func(ctx context.Context, userId string) (*[]models.DeliveryEvent, error)

type SaveAllDeliveryEvents func(ctx context.Context, data []models.DeliveryEvent) error

type SaveDeliveryEvent func(ctx context.Context, data models.DeliveryEvent) error

func FindAsStream(client *mongo.Client) FindDeliveriesAsStream {
	return func(ctx context.Context) (broadcast.Observer[*models.DeliveryEvent], error) {
		coll := client.Database(database).Collection(collection)

		channel, err := mongodb.FindAsStream[*models.DeliveryEvent](ctx, *coll)

		if err != nil {
			return nil, err
		}

		return channel, nil
	}
}

func FindByUserId(client *mongo.Client) FindDeliveriesByUserId {

	return func(ctx context.Context, userId string) (*[]models.DeliveryEvent, error) {
		coll := client.Database(database).Collection(collection)

		results, err := mongodb.FindByUserId[models.DeliveryEvent](ctx, *coll, userId)

		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func SaveAll(client *mongo.Client) SaveAllDeliveryEvents {
	return func(ctx context.Context, data []models.DeliveryEvent) error {
		coll := client.Database(database).Collection(collection)

		var items []interface{}

		for _, value := range data {
			items = append(items, value)
		}

		err := mongodb.SaveAll(ctx, *coll, items)
		return err
	}
}

func Save(client *mongo.Client) SaveDeliveryEvent {
	return func(ctx context.Context, data models.DeliveryEvent) error {
		coll := client.Database(database).Collection(collection)
		err := mongodb.Save(ctx, *coll, data)
		return err
	}
}
