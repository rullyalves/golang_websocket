package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateBlockDataParams struct {
	ID          string    `json:"id"`
	ReceiverID  string    `json:"receiverId"`
	ComplaintID *string   `json:"complaintId"`
	CreatedAt   time.Time `json:"createdAt"`
}

type SaveBlock func(context.Context, string, CreateBlockDataParams) error

func Save(driver *neo4j.DriverWithContext) SaveBlock {
	return func(ctx context.Context, userId string, data CreateBlockDataParams) error {

		params := map[string]any{
			"id":          data.ID,
			"senderId":    userId,
			"receiverId":  data.ReceiverID,
			"complaintId": data.ComplaintID,
			"createdAt":   data.CreatedAt,
		}

		query := `
			MATCH (sender :Profile{id: $senderId})  
			MATCH (receiver :Profile{id: $receiverId})
			WITH sender, receiver
			MERGE (sender)-[:SEND]->(b :Block)-[:RECEIVE]->(receiver)  
			ON CREATE SET
				b.id = COALESCE($id, b.id),
				b.createdAt = COALESCE($createdAt, b.createdAt)
				b.receiverId = $receiverId,
				b.complaintId = $complaintId
			WITH b
			MATCH (complaint: Complaint{id: $complaintId})  
			MERGE (b)-[:CONNECTED_TO]->(c)`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}
