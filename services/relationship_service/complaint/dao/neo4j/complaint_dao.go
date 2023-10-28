package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/relationship_service/complaint/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateComplaintDataParams struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	SenderID    string    `json:"sender_id"`
	ReceiverID  string    `json:"receiver_id"`
	CategoryID  string    `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type FindComplaintByBlockIdIn func(ctx context.Context, ids []string) ([]models.ComplaintView, error)

type SaveComplaint func(context.Context, CreateComplaintDataParams) error

func FindByBlockIdIn(driver *neo4j.DriverWithContext) FindComplaintByBlockIdIn {
	return func(ctx context.Context, ids []string) ([]models.ComplaintView, error) {
		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (b :Block) WHERE b.id IN ($ids) " +
				"MATCH (b)-->(c :Complaint) " +
				"RETURN b.id as ownerId, c"

		results, err := neo4jdb.ExecuteWithMapping[models.ComplaintView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveComplaint {
	return func(ctx context.Context, data CreateComplaintDataParams) error {

		params := map[string]any{
			"id":          data.ID,
			"description": data.Description,
			"senderId":    data.SenderID,
			"receiverId":  data.ReceiverID,
			"categoryId":  data.CategoryID,
			"createdAt":   data.CreatedAt,
		}

		query := `
			MATCH (sender: Profile{id: $senderId})  
			MATCH (receiver: Profile{id: $receiverId})  
			WITH sender, receiver  
			MERGE (sender)-[:SEND]->(c :Complaint)-[:RECEIVE]->(receiver)  
			ON CREATE SET
				c.id = COALESCE($id, c.id),
				c.createdAt = COALESCE($createdAt, c.createdAt)
				c.senderId = $senderId,
				c.receiverId = $receiverId,
				c.categoryId = $categoryId
			SET c.description = $description`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}
