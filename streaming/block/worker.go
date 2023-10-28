package block

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go_ws/shared/queue"
	"go_ws/shared/utils"
	"go_ws/streaming/block/dao"
	"go_ws/streaming/block/models"
)

func StartWorker(processQueue queue.ProcessMessage, saveAllEvents dao.SaveAllBlockEvents) func(queueUrl string) {
	worker := queue.StartWorker(processQueue, func(messages []types.Message, consume queue.MessageConsumer) error {
		data, err := utils.ParseMessagesTo[models.BlockEvent](messages)

		if err != nil {
			return err
		}

		err = saveAllEvents(context.Background(), data)

		if err != nil {
			return err
		}

		return consume(messages)
	})
	return worker
}
