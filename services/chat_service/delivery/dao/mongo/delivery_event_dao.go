package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	model "go_ws/services/chat_service/delivery/models"
	mongodb "go_ws/shared/database/mongo"
)

const (
	database   = "events"
	collection = "deliveries"
)

func FindAsStream(client *mongo.Client) func(ctx context.Context) (<-chan model.DeliveryEvent, error) {
	return func(ctx context.Context) (<-chan model.DeliveryEvent, error) {
		coll := client.Database(database).Collection(collection)

		channel, err := mongodb.FindAsStream[model.DeliveryEvent](ctx, *coll)

		if err != nil {
			return nil, err
		}

		return channel, nil
	}
}

func FindByUserId(client *mongo.Client) func(ctx context.Context, userId string) (*[]model.DeliveryEvent, error) {

	return func(ctx context.Context, userId string) (*[]model.DeliveryEvent, error) {
		coll := client.Database(database).Collection(collection)

		results, err := mongodb.FindByUserId[model.DeliveryEvent](ctx, *coll, userId)

		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(client *mongo.Client) func(ctx context.Context, data model.DeliveryEvent) error {
	return func(ctx context.Context, data model.DeliveryEvent) error {
		coll := client.Database(database).Collection(collection)
		err := mongodb.Save(ctx, *coll, data)
		return err
	}
}

func DeleteById(client *mongo.Client) func(ctx context.Context, id string) error {
	return func(ctx context.Context, id string) error {
		coll := client.Database(database).Collection(collection)
		return mongodb.DeleteById(ctx, *coll, id)
	}
}

func DeleteAll(ctx context.Context, client *mongo.Client) error {
	coll := client.Database(database).Collection(collection)
	return mongodb.DeleteAll(ctx, *coll)
}
