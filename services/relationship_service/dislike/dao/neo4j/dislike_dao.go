package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateDislikeDataParams struct {
	ID         string    `json:"id" validate:"required,uuid4"`
	ReceiverID string    `json:"receiverId" validate:"required,uuid4"`
	CreatedAt  time.Time `json:"createdAt" validate:"required"`
}

type SaveDislike func(context.Context, string, CreateDislikeDataParams) error

func Save(driver *neo4j.DriverWithContext) SaveDislike {
	return func(ctx context.Context, userId string, data CreateDislikeDataParams) error {

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
			MERGE (sender)-[:SEND]->(d :Dislike)-[:RECEIVE]->(receiver)
			ON CREATE SET
            	d.id = $id,
				d.createdAt = $createdAt,
				d.senderId = $senderId,
				d.receiverId = $receiverId`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}
