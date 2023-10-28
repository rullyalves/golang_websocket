package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go_ws/shared/database/mongo"
)

const (
	database   = "events"
	collection = "messages"
)

type DeleteMessageEvent func(ctx context.Context, id string) error

func DeleteById(client *mongo.Client) DeleteMessageEvent {
	return func(ctx context.Context, id string) error {
		coll := client.Database(database).Collection(collection)
		return mongodb.DeleteById(ctx, *coll, id)
	}
}
