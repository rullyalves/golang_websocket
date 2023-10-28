package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/relationship_service/like/models"
	"time"

	neo4jdb "go_ws/shared/database/neo4j"
)

type CreateLikeDataParams struct {
	ID         string    `json:"id"`
	ReceiverID string    `json:"receiver_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FindReceivedLikesByProfileId func(ctx context.Context, id string) (*[]models.LikeView, error)

type FindLikeBySenderAndReceiverId func(ctx context.Context, senderId string, receiverId string) (*models.LikeView, error)

type SaveLike func(context.Context, string, CreateLikeDataParams) error

type DeleteLikeById func(ctx context.Context, id string) error

func FindByProfileId(driver *neo4j.DriverWithContext) FindReceivedLikesByProfileId {
	return func(ctx context.Context, id string) (*[]models.LikeView, error) {
		params := map[string]any{
			"id": id,
		}

		query := `
		MATCH (p:Profile)-[send]->(l :Like)-[to]->(p2:Profile{id: $id}) 
		WHERE NOT (p)--(:Dislike|Match|Block)--(p2) 
		RETURN l`

		results, err := neo4jdb.ExecuteWithMapping[models.LikeView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindBySenderAndReceiverId(driver *neo4j.DriverWithContext) FindLikeBySenderAndReceiverId {
	return func(ctx context.Context, senderId string, receiverId string) (*models.LikeView, error) {

		params := map[string]any{
			"senderId":   senderId,
			"receiverId": receiverId,
		}

		query :=
			"MATCH (l:Like) " +
				"MATCH (:Profile{id: $senderId})-[s:SEND]->(l :Like)-[r:RECEIVE]->(:Profile{id: $receiverId}) " +
				"RETURN l"

		results, err := neo4jdb.ExecuteWithMapping[models.LikeView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		if len(results) == 0 {
			return nil, nil
		}

		return &results[0], nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveLike {
	return func(ctx context.Context, userId string, data CreateLikeDataParams) error {

		params := map[string]any{
			"id":         data.ID,
			"createdAt":  data.CreatedAt,
			"senderId":   userId,
			"receiverId": data.ReceiverID,
		}

		query := `
		MATCH (sender:Profile{id: $senderId}) 
		MATCH (receiver:Profile{id: $receiverId}) 
		WITH sender, receiver 
		MERGE (sender)-[:SEND]->(l :Like)-[:RECEIVE]->(receiver) 
		ON CREATE SET 
			l.id = $id, 
			l.createdAt = $createdAt, 
			l.senderId = $senderId, 
			l.receiverId = $receiverId`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteLikeById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		_, err := neo4jdb.ExecuteQuery(ctx, driver, "MATCH (i: Like{id: $id}) DETACH DELETE i", params)
		return err
	}
}
