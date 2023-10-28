package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/chat_service/attachment/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateAttachmentDataParams struct {
	ID          string    `json:"id"`
	ResourceURL string    `json:"resourceUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	MessageID   string    `json:"messageId"`
}

type FindAttachmentsByMessageId func(ctx context.Context, id string, limit int, cursor time.Time) ([]models.AttachmentView, error)

type FindAttachments func(ctx context.Context) ([]models.AttachmentView, error)

type SaveAttachment func(ctx context.Context, data CreateAttachmentDataParams) error

type DeleteAttachment func(ctx context.Context, id string) error

func FindByMessageId(driver *neo4j.DriverWithContext) FindAttachmentsByMessageId {
	return func(ctx context.Context, id string, limit int, cursor time.Time) ([]models.AttachmentView, error) {
		params := map[string]any{
			"id": id,
		}

		query :=
			"MATCH (Message{id: $id})-->(attachment :Attachment) " +
				"RETURN attachment"

		results, err := neo4jdb.ExecuteWithMapping[models.AttachmentView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func FindAll(driver *neo4j.DriverWithContext) FindAttachments {
	return func(ctx context.Context) ([]models.AttachmentView, error) {

		query := "MATCH (ch :Attachment) RETURN ch"

		results, err := neo4jdb.ExecuteWithMapping[models.AttachmentView](ctx, driver, query, nil)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveAttachment {
	return func(ctx context.Context, data CreateAttachmentDataParams) error {
		params := map[string]any{
			"id":          data.ID,
			"createdAt":   data.CreatedAt,
			"messageId":   data.MessageID,
			"resourceUrl": data.ResourceURL,
		}

		query := `
		MATCH (message :Message{id: $messageId}) 
		WITH message 
		MERGE (message)-[:HAS]->(at: Attachment) 
		ON CREATE SET 
			at.id = COALESCE($id, at.id),
			at.messageId = $messageId
			at.createdAt = COALESCE($createdAt, at.createdAt),
		 SET
			at.resourceUrl = COALESCE($resourceUrl, at.resourceUrl)`

		_, err := neo4jdb.ExecuteWithMapping[models.AttachmentView](ctx, driver, query, params)

		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteAttachment {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		query := "MATCH (ch :Attachment {id: $id}) DELETE ch"
		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}
