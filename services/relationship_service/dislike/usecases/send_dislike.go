package usecases

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	dislikeDao "go_ws/services/relationship_service/dislike/dao/neo4j"
	shared "go_ws/shared/models"
	"go_ws/shared/queue"
	streamingModels "go_ws/streaming/shared/models"
	"time"
)

type SendDislike func(context.Context, string, dislikeDao.CreateDislikeDataParams) error

func sendDislike(saveDislike dislikeDao.SaveDislike, sendToQueue queue.SendToQueue) SendDislike {

	return func(ctx context.Context, userId string, params dislikeDao.CreateDislikeDataParams) error {

		err := saveDislike(ctx, userId, params)

		if err != nil {
			return err
		}

		dislikeEvent := streamingModels.CreateEventPayload[any](
			//TODO: pegar dados do dislike
			nil,
			params.ReceiverID,
			time.Now().UTC(),
			shared.MessageTypeUpdate,
		)

		//TODO: pegar url da fila do env
		return sendToQueue(ctx, queue.DislikeStreamQueueUrl, []queue.MessageInput{dislikeEvent})
	}
}

func NewSendDislike(driver *neo4j.DriverWithContext, sqsClient *sqs.Client) SendDislike {
	saveDislike := dislikeDao.Save(driver)
	sendToQueue := queue.SendMessageToQueue(sqsClient)

	return sendDislike(saveDislike, sendToQueue)
}
