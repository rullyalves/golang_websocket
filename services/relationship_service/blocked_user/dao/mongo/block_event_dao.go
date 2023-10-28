package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go_ws/shared/database/mongo"
)

const (
	database   = "events"
	collection = "blocked_users"
)

type DeleteBlockEventById func(ctx context.Context, id string) error

func DeleteById(client *mongo.Client) DeleteBlockEventById {
	return func(ctx context.Context, id string) error {
		coll := client.Database(database).Collection(collection)
		return mongodb.DeleteById(ctx, *coll, id)
	}
}
