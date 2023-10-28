package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_ws/streaming/shared/broadcast"
	"log"
)

type changeEvent[T any] struct {
	FullDocument T
}

func FindAsStream[T interface{}](ctx context.Context, collection mongo.Collection) (broadcast.Observer[T], error) {

	pipeline := mongo.Pipeline{
		{{
			"$match", bson.D{{
				"operationType", bson.D{{
					"$in", bson.A{
						"insert", "update", "replace",
					},
				}},
			}},
		}},
	}

	stream, err := collection.Watch(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	br := broadcast.New[T](5)

	go func(ctx context.Context) {
		defer func() {
			if err := stream.Close(ctx); err != nil {
				log.Fatal(err)
			}
		}()

		for stream.Next(ctx) {
			var item changeEvent[T]

			if err := stream.Decode(&item); err != nil {
				log.Fatal(err)
			}

			br.Publish(item.FullDocument)
		}
	}(ctx)

	return br, nil
}

func FindByUserId[T interface{}](ctx context.Context, collection mongo.Collection, userId string) (*[]T, error) {

	pipeline := bson.D{{Key: "userId", Value: userId}}

	cursor, err := collection.Find(ctx, pipeline)

	defer func(ctx context.Context) {
		if err := cursor.Close(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	if err != nil {
		return nil, err
	}

	var results []T

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func Save(ctx context.Context, collection mongo.Collection, document interface{}) error {
	_, err := collection.InsertOne(ctx, document)
	return err
}

func SaveAll(ctx context.Context, collection mongo.Collection, documents []interface{}) error {
	opts := options.InsertMany().SetOrdered(false)
	_, err := collection.InsertMany(ctx, documents, opts)
	return err
}

func DeleteById(ctx context.Context, collection mongo.Collection, id string) error {
	_, err := collection.DeleteOne(ctx, bson.D{{"id", id}})
	return err
}

func DeleteAll(ctx context.Context, collection mongo.Collection) error {
	_, err := collection.DeleteMany(ctx, bson.D{})
	return err
}
