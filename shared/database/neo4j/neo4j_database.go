package neo4j

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type transactionKey string

type TransactionRunner[T any] func(ctx context.Context) (*T, error)

const (
	transaction transactionKey = "neo4jTransaction"
)

func NewDriver(uri, username, password string) (*neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))

	if err != nil {
		return nil, err
	}

	if connectionErr := driver.VerifyConnectivity(context.Background()); connectionErr != nil {
		return nil, connectionErr
	}

	fmt.Println("Pinged your deployment. You successfully connected to Neo4j!")

	return &driver, nil
}

func WithNewSession(ctx context.Context, driver *neo4j.DriverWithContext) neo4j.SessionWithContext {
	session := (*driver).NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	return session
}

func WithTransaction[T any](driver *neo4j.DriverWithContext) func(ctx context.Context, runTransaction TransactionRunner[T]) (*T, error) {

	return func(ctx context.Context, runTransaction TransactionRunner[T]) (*T, error) {
		session := WithNewSession(ctx, driver)

		tx, beginTxErr := session.BeginTransaction(ctx)

		if beginTxErr != nil {
			return nil, beginTxErr
		}

		defer func() {
			sessionErr := session.Close(ctx)
			txErr := tx.Close(ctx)

			if sessionErr != nil {
				panic(sessionErr)
			}

			if txErr != nil {
				panic(txErr)
			}
		}()

		ctx = context.WithValue(ctx, transaction, tx)

		result, executionErr := runTransaction(ctx)

		if executionErr != nil {
			rollbackErr := tx.Rollback(ctx)

			if rollbackErr != nil {
				return nil, rollbackErr
			}

			return nil, executionErr
		}

		commitErr := tx.Commit(ctx)

		if commitErr != nil {
			return nil, commitErr
		}

		return result, nil
	}
}

func ExecuteQuery(ctx context.Context, driver *neo4j.DriverWithContext, cypher string, params map[string]any) (neo4j.ResultWithContext, error) {
	tx, _ := ctx.Value(transaction).(neo4j.ExplicitTransaction)

	if tx != nil {
		return tx.Run(ctx, cypher, params)
	}

	newSession := WithNewSession(ctx, driver)

	return newSession.Run(ctx, cypher, params)
}

func mapField(node any) any {

	var result any

	switch properties := node.(type) {
	case dbtype.Node:
		var results = make(map[string]any)

		for key, property := range properties.Props {
			results[key] = mapField(property)
		}
		result = results
	case map[string]any:
		var results = make(map[string]any)

		for key, property := range properties {
			results[key] = mapField(property)
		}
		result = results
	case dbtype.LocalDateTime:
		result = properties.Time()
	case []interface{}:
		var results = make([]any, len(properties))

		for i, property := range properties {
			results[i] = mapField(property)
		}
		result = results
	default:
		result = properties
	}

	return result
}

func ExecuteWithMapping[T interface{}](
	ctx context.Context,
	driver *neo4j.DriverWithContext,
	query string,
	params map[string]any,
) ([]T, error) {

	result, err := ExecuteQuery(ctx, driver, query, params)

	if err != nil {
		return nil, err
	}

	records, resultErr := result.Collect(ctx)

	if resultErr != nil {
		return nil, resultErr
	}

	var results []T = make([]T, 0)

	for _, record := range records {

		properties := mapField(record.Values[0])

		var newData T
		err := mapstructure.Decode(properties, &newData)

		if err != nil {
			return nil, err
		}

		results = append(results, newData)
	}

	return results, nil
}
