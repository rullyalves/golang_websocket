package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go_ws/shared/database/mongo"
	"go_ws/streaming/like/models"
	"go_ws/streaming/shared/broadcast"
)

const (
	database   = "events"
	collection = "likes"
)

type FindLikeEventsAsStream func(ctx context.Context) (broadcast.Observer[*models.LikeEvent], error)

type FindLikeEventsByUserId func(ctx context.Context, userId string) (*[]models.LikeEvent, error)

type SaveAllLikeEvents func(ctx context.Context, data []models.LikeEvent) error

type SaveLikeEvent func(ctx context.Context, data models.LikeEvent) error

func FindAsStream(client *mongo.Client) FindLikeEventsAsStream {
	return func(ctx context.Context) (broadcast.Observer[*models.LikeEvent], error) {
		coll := client.Database(database).Collection(collection)

		channel, err := mongodb.FindAsStream[*models.LikeEvent](ctx, *coll)

		if err != nil {
			return nil, err
		}

		return channel, nil
	}
}

func FindByUserId(client *mongo.Client) FindLikeEventsByUserId {

	return func(ctx context.Context, userId string) (*[]models.LikeEvent, error) {
		coll := client.Database(database).Collection(collection)

		results, err := mongodb.FindByUserId[models.LikeEvent](ctx, *coll, userId)

		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func SaveAll(client *mongo.Client) SaveAllLikeEvents {
	return func(ctx context.Context, data []models.LikeEvent) error {
		coll := client.Database(database).Collection(collection)

		var items []interface{}

		for _, value := range data {
			items = append(items, value)
		}

		err := mongodb.SaveAll(ctx, *coll, items)
		return err
	}
}

func Save(client *mongo.Client) SaveLikeEvent {
	return func(ctx context.Context, data models.LikeEvent) error {
		coll := client.Database(database).Collection(collection)
		err := mongodb.Save(ctx, *coll, data)
		return err
	}
}
