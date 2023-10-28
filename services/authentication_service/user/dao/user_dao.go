package dao

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/authentication_service/user/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateUserDataParams struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

type FindUserByIdIn func(ctx context.Context, ids []string) ([]models.UserView, error)

type FindUserByUsername func(ctx context.Context, username string) (*models.UserView, error)

type FindUsers func(ctx context.Context) ([]models.UserView, error)

type SaveUser func(context.Context, CreateUserDataParams) (*models.UserView, error)

type DeleteUserById func(ctx context.Context, id string) error

func FindByUserIdIn(driver *neo4j.DriverWithContext) FindUserByIdIn {
	return func(ctx context.Context, ids []string) ([]models.UserView, error) {
		params := map[string]any{
			"ids": ids,
		}

		query := "MATCH (p: User) WHERE p.id IN ($ids) "

		results, err := neo4jdb.ExecuteWithMapping[models.UserView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func FindByUsername(driver *neo4j.DriverWithContext) FindUserByUsername {
	return func(ctx context.Context, username string) (*models.UserView, error) {
		params := map[string]any{
			"username": username,
		}

		query := "MATCH (p: User{username: $username}) RETURN p "

		results, err := neo4jdb.ExecuteWithMapping[models.UserView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		if len(results) == 0 {
			return nil, nil
		}

		return &results[0], nil
	}
}

func FindAll(driver *neo4j.DriverWithContext) FindUsers {
	return func(ctx context.Context) ([]models.UserView, error) {

		query := "MATCH (n: User) RETURN n"

		results, err := neo4jdb.ExecuteWithMapping[models.UserView](ctx, driver, query, nil)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveUser {
	return func(ctx context.Context, data CreateUserDataParams) (*models.UserView, error) {

		params := map[string]any{
			"id":        data.ID,
			"username":  data.Username,
			"createdAt": data.CreatedAt,
		}

		query := `
		MERGE (n: User{username: $username}) 
		ON CREATE SET 
			n.id = COALESCE($id, n.id),
			n.username = $username,
 			n.createdAt = COALESCE($createdAt, n.createdAt) 
		RETURN n`

		result, err := neo4jdb.ExecuteWithMapping[models.UserView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		if len(result) == 0 {
			return nil, nil
		}

		return &result[0], nil
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteUserById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}

		query := "MATCH (n: User{id: $id}) DETACH DELETE n"

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}
