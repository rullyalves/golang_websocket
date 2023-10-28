package mongo

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

type FindChatsAsStream func(ctx context.Context) (broadcast.Observer[models.ChatEvent], error)

type FindChatsByUserId func(ctx context.Context, userId string) (*[]models.ChatEvent, error)

type SaveChatEvent func(ctx context.Context, data models.ChatEvent) error

type DeleteChatEvent func(ctx context.Context, id string) error

type DeleteAllChatEvents func(ctx context.Context) error

func FindAsStream(client *mongo.Client) FindChatsAsStream {
	return func(ctx context.Context) (broadcast.Observer[models.ChatEvent], error) {
		coll := client.Database(database).Collection(collection)

		channel, err := mongodb.FindAsStream[models.ChatEvent](ctx, *coll)

		if err != nil {
			return nil, err
		}

		return channel, nil
	}
}

func FindByUserId(client *mongo.Client) FindChatsByUserId {

	return func(ctx context.Context, userId string) (*[]models.ChatEvent, error) {
		coll := client.Database(database).Collection(collection)

		results, err := mongodb.FindByUserId[models.ChatEvent](ctx, *coll, userId)

		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(client *mongo.Client) SaveChatEvent {
	return func(ctx context.Context, data models.ChatEvent) error {
		coll := client.Database(database).Collection(collection)
		err := mongodb.Save(ctx, *coll, data)
		return err
	}
}

func DeleteById(client *mongo.Client) DeleteChatEvent {
	return func(ctx context.Context, id string) error {
		coll := client.Database(database).Collection(collection)
		return mongodb.DeleteById(ctx, *coll, id)
	}
}

func DeleteAll(client *mongo.Client) DeleteAllChatEvents {
	return func(ctx context.Context) error {
		coll := client.Database(database).Collection(collection)
		return mongodb.DeleteAll(ctx, *coll)
	}
}
