package usecases

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	deliveryDao "go_ws/services/chat_service/delivery/dao/neo4j"
	"go_ws/services/chat_service/delivery/models"
	shared "go_ws/shared/models"
	"go_ws/shared/queue"
	streamingModels "go_ws/streaming/shared/models"
	"time"
)

type CreateDeliveryParams struct {
	ID        string                `json:"id" validate:"required,uuid4"`
	CreatedAt time.Time             `json:"createdAt" validate:"required"`
	Status    models.DeliveryStatus `json:"status" validate:"required,oneof=sent delivered read played"`
	MessageID string                `json:"messageId" validate:"required,uuid4"`
	TargetID  string                `json:"targetId" validate:"required,uuid4"`
}

type ConfirmMessageDelivery func(context.Context, CreateDeliveryParams) error

func NewConfirmMessageDelivery(neo4jDriver *neo4j.DriverWithContext, sqsClient *sqs.Client) ConfirmMessageDelivery {
	saveDelivery := deliveryDao.Save(neo4jDriver)
	sendToQueue := queue.SendMessageToQueue(sqsClient)
	return confirmDelivery(saveDelivery, sendToQueue)
}

func confirmDelivery(saveDelivery deliveryDao.SaveDelivery, sendToQueue queue.SendToQueue) ConfirmMessageDelivery {
	return func(ctx context.Context, params CreateDeliveryParams) error {

		err := saveDelivery(ctx, deliveryDao.CreateDeliveryDataParams{
			ID:        params.ID,
			CreatedAt: params.CreatedAt,
			Status:    params.Status,
			MessageID: params.MessageID,
			TargetID:  params.TargetID,
		})

		if err != nil {
			return err
		}

		return sendEvents(sendToQueue)(ctx, params)
	}
}

func sendEvents(sendToQueue queue.SendToQueue) func(ctx context.Context, data CreateDeliveryParams) error {
	return func(ctx context.Context, params CreateDeliveryParams) error {

		deliveryEvent := streamingModels.CreateEventPayload[models.DeliveryView](
			models.DeliveryView{
				ID:        params.ID,
				CreatedAt: params.CreatedAt,
				Status:    params.Status,
				TargetID:  params.TargetID,
				MessageID: params.MessageID,
			},
			params.TargetID,
			time.Now().UTC(),
			shared.MessageTypeUpdate,
		)

		return sendToQueue(ctx, queue.DeliveryStreamQueueUrl, []queue.MessageInput{deliveryEvent})
	}
}
