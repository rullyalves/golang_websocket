package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go_ws/shared/database/mongo"
	"go_ws/streaming/chat/models"
	"go_ws/streaming/shared/broadcast"
)

const (
	database   = "events"
	collection = "chats"
)

type FindChatEventsAsStream func(ctx context.Context) (broadcast.Observer[*models.ChatEvent], error)

type FindChatsEventsByUserId func(ctx context.Context, userId string) (*[]models.ChatEvent, error)

type SaveChatEvent func(ctx context.Context, data models.ChatEvent) error

type SaveAllChatEvents func(ctx context.Context, data []models.ChatEvent) error

func FindAsStream(client *mongo.Client) FindChatEventsAsStream {
	return func(ctx context.Context) (broadcast.Observer[*models.ChatEvent], error) {
		coll := client.Database(database).Collection(collection)

		channel, err := mongodb.FindAsStream[*models.ChatEvent](ctx, *coll)

		if err != nil {
			return nil, err
		}

		return channel, nil
	}
}

func FindByUserId(client *mongo.Client) FindChatsEventsByUserId {

	return func(ctx context.Context, userId string) (*[]models.ChatEvent, error) {
		coll := client.Database(database).Collection(collection)

		results, err := mongodb.FindByUserId[models.ChatEvent](ctx, *coll, userId)

		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func SaveAll(client *mongo.Client) SaveAllChatEvents {
	return func(ctx context.Context, data []models.ChatEvent) error {
		coll := client.Database(database).Collection(collection)

		var items []interface{}

		for _, value := range data {
			items = append(items, value)
		}

		err := mongodb.SaveAll(ctx, *coll, items)
		return err
	}
}

func Save(client *mongo.Client) SaveChatEvent {
	return func(ctx context.Context, data models.ChatEvent) error {
		coll := client.Database(database).Collection(collection)
		err := mongodb.Save(ctx, *coll, data)

		return err
	}
}
