package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionRunner func(ctx context.Context) error

func NewDriver(uri string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}

func WithNewSession(ctx context.Context, client *mongo.Client) (mongo.Session, error) {
	session, err := client.StartSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func WithTransaction(ctx context.Context, driver *mongo.Client, runTransaction TransactionRunner) error {

	session, sessionErr := WithNewSession(ctx, driver)

	if sessionErr != nil {
		return sessionErr
	}

	beginTxErr := session.StartTransaction()

	if beginTxErr != nil {
		return beginTxErr
	}

	defer func() {
		session.EndSession(ctx)
	}()

	executionErr := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		return runTransaction(sc)
	})

	if executionErr != nil {
		rollbackErr := session.AbortTransaction(ctx)

		if rollbackErr != nil {
			return rollbackErr
		}

		return executionErr
	}

	commitErr := session.CommitTransaction(ctx)

	if commitErr != nil {
		return commitErr
	}

	return nil
}
