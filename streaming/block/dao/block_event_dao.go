package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go_ws/shared/database/mongo"
	"go_ws/streaming/block/models"
	"go_ws/streaming/shared/broadcast"
)

const (
	database   = "events"
	collection = "blocks"
)

type FindBlockEventsAsStream func(ctx context.Context) (broadcast.Observer[*models.BlockEvent], error)

type FindBlockEventsByUserId func(ctx context.Context, userId string) (*[]models.BlockEvent, error)

type SaveBlockEvent func(ctx context.Context, data models.BlockEvent) error

type SaveAllBlockEvents func(ctx context.Context, data []models.BlockEvent) error

func FindAsStream(client *mongo.Client) FindBlockEventsAsStream {
	return func(ctx context.Context) (broadcast.Observer[*models.BlockEvent], error) {
		coll := client.Database(database).Collection(collection)

		channel, err := mongodb.FindAsStream[*models.BlockEvent](ctx, *coll)

		if err != nil {
			return nil, err
		}

		return channel, nil
	}
}

func FindByUserId(client *mongo.Client) FindBlockEventsByUserId {

	return func(ctx context.Context, userId string) (*[]models.BlockEvent, error) {
		coll := client.Database(database).Collection(collection)

		results, err := mongodb.FindByUserId[models.BlockEvent](ctx, *coll, userId)

		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func SaveAll(client *mongo.Client) SaveAllBlockEvents {
	return func(ctx context.Context, data []models.BlockEvent) error {
		coll := client.Database(database).Collection(collection)

		var items []interface{}

		for _, value := range data {
			items = append(items, value)
		}

		err := mongodb.SaveAll(ctx, *coll, items)
		return err
	}
}

func Save(client *mongo.Client) SaveBlockEvent {
	return func(ctx context.Context, data models.BlockEvent) error {
		coll := client.Database(database).Collection(collection)
		err := mongodb.Save(ctx, *coll, data)
		return err
	}
}
