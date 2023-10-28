package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go_ws/shared/database/mongo"
	"go_ws/streaming/message/models"
	"go_ws/streaming/shared/broadcast"
)

const (
	database   = "events"
	collection = "messages"
)

type FindMessageEventsAsStream func(ctx context.Context) (broadcast.Observer[*models.MessageEvent], error)

type FindMessageEventsByUserId func(ctx context.Context, userId string) (*[]models.MessageEvent, error)

type SaveAllMessageEvents func(ctx context.Context, data []models.MessageEvent) error

type SaveMessageEvent func(ctx context.Context, data models.MessageEvent) error

func FindAsStream(client *mongo.Client) FindMessageEventsAsStream {
	return func(ctx context.Context) (broadcast.Observer[*models.MessageEvent], error) {
		coll := client.Database(database).Collection(collection)

		channel, err := mongodb.FindAsStream[*models.MessageEvent](ctx, *coll)

		if err != nil {
			return nil, err
		}

		return channel, nil
	}
}

func FindByUserId(client *mongo.Client) FindMessageEventsByUserId {

	return func(ctx context.Context, userId string) (*[]models.MessageEvent, error) {
		coll := client.Database(database).Collection(collection)

		results, err := mongodb.FindByUserId[models.MessageEvent](ctx, *coll, userId)

		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func SaveAll(client *mongo.Client) SaveAllMessageEvents {
	return func(ctx context.Context, data []models.MessageEvent) error {
		coll := client.Database(database).Collection(collection)

		var items []interface{}

		for _, value := range data {
			items = append(items, value)
		}

		err := mongodb.SaveAll(ctx, *coll, items)
		return err
	}
}

func Save(client *mongo.Client) SaveMessageEvent {
	return func(ctx context.Context, data models.MessageEvent) error {
		coll := client.Database(database).Collection(collection)
		err := mongodb.Save(ctx, *coll, data)
		return err
	}
}

func DeleteAll(client *mongo.Client) func(context.Context) error {
	return func(ctx context.Context) error {
		coll := client.Database(database).Collection(collection)
		err := mongodb.DeleteAll(ctx, *coll)
		return err
	}
}
