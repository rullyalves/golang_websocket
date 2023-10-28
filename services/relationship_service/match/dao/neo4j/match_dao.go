package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/relationship_service/match/models"
	"time"

	neo4jdb "go_ws/shared/database/neo4j"
)

type CreateMatchDataParams struct {
	ID        string    `json:"id"`
	LikeIds   []string  `json:"like_ids"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type FindMatchByProfileIdIn func(ctx context.Context, ids []string) ([]models.MatchOwnerView, error)

type SaveMatch func(context.Context, CreateMatchDataParams) error

type DeleteMatchById func(ctx context.Context, id string) error

func FindProfileIdIn(driver *neo4j.DriverWithContext) FindMatchByProfileIdIn {
	return func(ctx context.Context, ids []string) ([]models.MatchOwnerView, error) {
		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (p: Profile) WHERE p.id IN ($ids) " +
				"MATCH (p)--(:Like)--(mh: Match) " +
				"RETURN {ownerId: p.id, data: mh}"

		results, err := neo4jdb.ExecuteWithMapping[models.MatchOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveMatch {
	return func(ctx context.Context, data CreateMatchDataParams) error {

		params := map[string]any{
			"id":        data.ID,
			"likeIds":   data.LikeIds,
			"isActive":  data.IsActive,
			"createdAt": data.CreatedAt,
		}

		query := `
			MATCH (like: Like) WHERE like.id IN ($likeIds)  
			WITH like  
			MERGE (mh: Match{id: $id})  
			ON CREATE SET
				mh.likeIds = $likeIds,
				mh.createdAt = COALESCE($createdAt, mh.createdAt)
			SET
				mh.isActive = $isActive
			MERGE (like)-[:COMPOSE]->(mh)`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteMatchById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		_, err := neo4jdb.ExecuteQuery(ctx, driver, "MATCH (i: Match{id: $id}) DETACH DELETE i", params)
		return err
	}
}
