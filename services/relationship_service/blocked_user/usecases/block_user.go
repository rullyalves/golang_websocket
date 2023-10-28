package usecases

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	blockDao "go_ws/services/relationship_service/blocked_user/dao/neo4j"
	shared "go_ws/shared/models"
	"go_ws/shared/queue"
	streamingModels "go_ws/streaming/shared/models"
	"time"
)

type BlockUser func(context.Context, string, blockDao.CreateBlockDataParams) error

func blockUser(saveBlock blockDao.SaveBlock, sendToQueue queue.SendToQueue) BlockUser {

	return func(ctx context.Context, userId string, params blockDao.CreateBlockDataParams) error {

		err := saveBlock(ctx, userId, params)

		if err != nil {
			return err
		}

		blockEvent := streamingModels.CreateEventPayload[any](
			//TODO: pegar dados do bloqueio
			nil,
			params.ReceiverID,
			time.Now().UTC(),
			shared.MessageTypeUpdate,
		)

		return sendToQueue(ctx, queue.BlockStreamQueueUrl, []queue.MessageInput{blockEvent})
	}
}

func NewBlockUser(driver *neo4j.DriverWithContext, sqsClient *sqs.Client) BlockUser {
	saveBlock := blockDao.Save(driver)
	sendToQueue := queue.SendMessageToQueue(sqsClient)

	return blockUser(saveBlock, sendToQueue)
}
