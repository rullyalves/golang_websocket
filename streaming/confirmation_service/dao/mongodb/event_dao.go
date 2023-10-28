package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go_ws/shared/database/mongo"
)

const (
	database = "events"
)

type Event string

const (
	blocks     Event = "blocks"
	chats      Event = "chats"
	deliveries Event = "deliveries"
	messages   Event = "messages"
	profiles   Event = "profiles"
	dislikes   Event = "dislikes"
	likes      Event = "likes"
)

type DeleteEventById func(ctx context.Context, id string, event Event) error

func DeleteById(client *mongo.Client) DeleteEventById {
	return func(ctx context.Context, id string, event Event) error {
		coll := client.Database(database).Collection(string(event))
		return mongodb.DeleteById(ctx, *coll, id)
	}
}
