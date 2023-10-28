package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/profile_service/feedback/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateFeedbackDataParams struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type FindFeedbackByProfileIdIn func(ctx context.Context, ids []string) ([]models.FeedbackOwnerView, error)

type SaveFeedback func(context.Context, string, CreateFeedbackDataParams) error

type DeleteFeedbackById func(ctx context.Context, id string) error

func FindProfileIdIn(driver *neo4j.DriverWithContext) FindFeedbackByProfileIdIn {
	return func(ctx context.Context, ids []string) ([]models.FeedbackOwnerView, error) {
		params := map[string]any{
			"ids": ids,
		}

		query := "MATCH (p:Profile)-[h]->(a:Feedback) WHERE p.id IN ($ids) RETURN {ownerId: p.id, data: a}"

		results, err := neo4jdb.ExecuteWithMapping[models.FeedbackOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveFeedback {
	return func(ctx context.Context, userId string, data CreateFeedbackDataParams) error {
		params := map[string]any{
			"id":          data.ID,
			"description": data.Description,
			"createdAt":   data.CreatedAt,
			"profileId":   userId,
		}

		//TODO: verificar se essa query possibilita mandar mais de um feedback por usuário

		query := `
		MATCH (sender: Profile{id: $profileId}) 
		WITH sender 
		MERGE (sender)-[:SEND]->(fb: Feedback{id: $id}) 
		ON CREATE SET
			fb.createdAt = COALESCE($createdAt, fb.createdAt)
		SET 
			fb.description = $description`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteFeedbackById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		_, err := neo4jdb.ExecuteQuery(ctx, driver, "MATCH (i: Feedback{id: $id}) DETACH DELETE i", params)
		return err
	}
}
