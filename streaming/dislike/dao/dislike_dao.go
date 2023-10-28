package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go_ws/shared/database/mongo"
	"go_ws/streaming/dislike/models"
	"go_ws/streaming/shared/broadcast"
)

const (
	database   = "events"
	collection = "dislikes"
)

type FindDislikeEventsAsStream func(ctx context.Context) (broadcast.Observer[*models.DislikeEvent], error)

type FindDislikeEventsByUserId func(ctx context.Context, userId string) (*[]models.DislikeEvent, error)

type SaveAllDislikeEvents func(ctx context.Context, data []models.DislikeEvent) error

type SaveDislikeEvent func(ctx context.Context, data models.DislikeEvent) error

func FindAsStream(client *mongo.Client) FindDislikeEventsAsStream {
	return func(ctx context.Context) (broadcast.Observer[*models.DislikeEvent], error) {
		coll := client.Database(database).Collection(collection)

		channel, err := mongodb.FindAsStream[*models.DislikeEvent](ctx, *coll)

		if err != nil {
			return nil, err
		}

		return channel, nil
	}
}

func FindByUserId(client *mongo.Client) FindDislikeEventsByUserId {

	return func(ctx context.Context, userId string) (*[]models.DislikeEvent, error) {
		coll := client.Database(database).Collection(collection)

		results, err := mongodb.FindByUserId[models.DislikeEvent](ctx, *coll, userId)

		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func SaveAll(client *mongo.Client) SaveAllDislikeEvents {
	return func(ctx context.Context, data []models.DislikeEvent) error {
		coll := client.Database(database).Collection(collection)

		var items []interface{}

		for _, value := range data {
			items = append(items, value)
		}

		err := mongodb.SaveAll(ctx, *coll, items)
		return err
	}
}

func Save(client *mongo.Client) SaveDislikeEvent {
	return func(ctx context.Context, data models.DislikeEvent) error {
		coll := client.Database(database).Collection(collection)
		err := mongodb.Save(ctx, *coll, data)
		return err
	}
}
