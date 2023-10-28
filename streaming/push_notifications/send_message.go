package messaging

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go_ws/shared/queue"
	"go_ws/shared/utils"
)

func SendNotificationFromQueue(sendNotification SendNotification) queue.MessageProcessor {
	return func(messages []types.Message, consume queue.MessageConsumer) error {
		ctx := context.Background()
		data, err := utils.ParseMessagesTo[*NotificationParams](messages)

		if err != nil {
			return err
		}
		err = sendNotification(ctx, data)

		if err != nil {
			return err
		}

		return consume(messages)
	}
}
