package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/chat_service/delivery/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type ChatLastMessagesResult struct {
	MessageID  string                `json:"messageId"`
	Deliveries []models.DeliveryView `json:"deliveries"`
}

type CreateDeliveryDataParams struct {
	ID        string                `json:"id"`
	CreatedAt time.Time             `json:"createdAt"`
	Status    models.DeliveryStatus `json:"status"`
	MessageID string                `json:"messageId"`
	TargetID  string                `json:"targetId"`
}

type FindDeliveriesByMessageIdIn func(ctx context.Context, ids []string) (map[string][]models.DeliveryView, error)
type FindDeliveryByMessageAndReceiverIds func(ctx context.Context, messageId string, receiverId string) ([]models.DeliveryView, error)
type SaveDelivery func(ctx context.Context, data CreateDeliveryDataParams) error
type DeleteDeliveryById func(ctx context.Context, id string) error

func FindByMessageIdIn(driver *neo4j.DriverWithContext) FindDeliveriesByMessageIdIn {
	return func(ctx context.Context, ids []string) (map[string][]models.DeliveryView, error) {

		params := map[string]any{
			"ids": ids,
		}

		query := `
		MATCH (message: Message) WHERE message.id IN ($ids)
		CALL {
			WITH message
			MATCH (delivery: Delivery) WHERE (message)-->(delivery) RETURN delivery ORDER BY delivery.createdAt ASC
		}
		WITH message.id as ownerId, collect(delivery) as deliveries
		RETURN {messageId: ownerId, deliveries: deliveries}`

		result, err := neo4jdb.ExecuteWithMapping[ChatLastMessagesResult](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		mapResult := make(map[string][]models.DeliveryView)

		for _, item := range result {
			mapResult[item.MessageID] = item.Deliveries
		}

		return mapResult, nil
	}
}

func FindByMessageIdAndReceiverId(driver *neo4j.DriverWithContext) FindDeliveryByMessageAndReceiverIds {
	return func(ctx context.Context, messageId string, receiverId string) ([]models.DeliveryView, error) {
		params := map[string]any{
			"messageId":  messageId,
			"receiverId": receiverId,
		}

		query :=
			"MATCH (Message{id: $messageId})-->(dl: Delivery)<--(Profile{id: $receiverId}) " +
				"RETURN dl"

		results, err := neo4jdb.ExecuteWithMapping[models.DeliveryView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveDelivery {
	return func(ctx context.Context, data CreateDeliveryDataParams) error {
		params := map[string]any{
			"id":        data.ID,
			"createdAt": data.CreatedAt,
			"status":    data.Status,
			"messageId": data.MessageID,
			"targetId":  data.TargetID,
		}

		query := `
		MATCH (ms :Message{id: $messageId}) 
		MATCH (p :Profile{id: $targetId}) 
		WITH ms, p 
		MERGE (ms)-[:FROM]->(dl: Delivery{id: $id})<-[:TO]-(p) 
		ON CREATE SET 
			dl.targetId = $targetId,
			dl.messageId = $messageId,
			dl.createdAt = COALESCE($createdAt, dl.createdAt)
		SET
			dl.status = $status`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteDeliveryById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		query := "MATCH (dl :Delivery{id: $id}) DELETE dl"
		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}
