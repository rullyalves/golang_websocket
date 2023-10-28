package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go_ws/shared/database/mongo"
	"go_ws/streaming/profile/models"
	"go_ws/streaming/shared/broadcast"
)

const (
	database   = "events"
	collection = "profiles"
)

type FindProfileEventsAsStream func(ctx context.Context) (broadcast.Observer[*models.ProfileEvent], error)

type FindProfileEventsByUserId func(ctx context.Context, userId string) (*[]models.ProfileEvent, error)

type SaveAllProfileEvents func(ctx context.Context, data []models.ProfileEvent) error

type SaveProfileEvent func(ctx context.Context, data models.ProfileEvent) error

func FindAsStream(client *mongo.Client) FindProfileEventsAsStream {
	return func(ctx context.Context) (broadcast.Observer[*models.ProfileEvent], error) {
		coll := client.Database(database).Collection(collection)

		channel, err := mongodb.FindAsStream[*models.ProfileEvent](ctx, *coll)

		if err != nil {
			return nil, err
		}

		return channel, nil
	}
}

func FindByUserId(client *mongo.Client) FindProfileEventsByUserId {

	return func(ctx context.Context, userId string) (*[]models.ProfileEvent, error) {
		coll := client.Database(database).Collection(collection)

		results, err := mongodb.FindByUserId[models.ProfileEvent](ctx, *coll, userId)

		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func SaveAll(client *mongo.Client) SaveAllProfileEvents {
	return func(ctx context.Context, data []models.ProfileEvent) error {
		coll := client.Database(database).Collection(collection)

		var items []interface{}

		for _, value := range data {
			items = append(items, value)
		}

		err := mongodb.SaveAll(ctx, *coll, items)
		return err
	}
}

func Save(client *mongo.Client) SaveProfileEvent {
	return func(ctx context.Context, data models.ProfileEvent) error {
		coll := client.Database(database).Collection(collection)
		err := mongodb.Save(ctx, *coll, data)
		return err
	}
}
