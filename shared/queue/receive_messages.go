package queue

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log"
)

type MessageConsumer func(message []types.Message) error

type MessageProcessor func(messages []types.Message, consume MessageConsumer) error

type ProcessMessage func(queueName string, size int32, process MessageProcessor) error

type ProcessChunk func(ctx context.Context, chunk []types.Message, consume MessageConsumer)

func ProcessQueue(sqsClient *sqs.Client) ProcessMessage {
	return func(queueName string, size int32, process MessageProcessor) error {

		ctx := context.Background()

		msgConfig := sqs.ReceiveMessageInput{
			MaxNumberOfMessages: size,
			QueueUrl:            &queueName,
			VisibilityTimeout:   20,
			WaitTimeSeconds:     20,
		}

		result, err := sqsClient.ReceiveMessage(ctx, &msgConfig)

		if err != nil {
			log.Printf("unable to receive message, %v", err)
			return err
		}

		messages := result.Messages

		if len(messages) == 0 {
			return nil
		}

		var changes []types.ChangeMessageVisibilityBatchRequestEntry

		for _, message := range result.Messages {
			input := types.ChangeMessageVisibilityBatchRequestEntry{
				Id:                message.MessageId,
				VisibilityTimeout: 20,
				ReceiptHandle:     message.ReceiptHandle,
			}
			changes = append(changes, input)
		}

		_, err = sqsClient.ChangeMessageVisibilityBatch(ctx,
			&sqs.ChangeMessageVisibilityBatchInput{
				QueueUrl: &queueName,
				Entries:  changes,
			},
		)

		if err != nil {
			return err
		}

		err = process(messages, func(messages []types.Message) error {

			var changes []types.DeleteMessageBatchRequestEntry

			for _, message := range messages {
				input := types.DeleteMessageBatchRequestEntry{
					Id:            message.MessageId,
					ReceiptHandle: message.ReceiptHandle,
				}
				changes = append(changes, input)
			}

			input := sqs.DeleteMessageBatchInput{
				QueueUrl: &queueName,
				Entries:  changes,
			}
			_, err := sqsClient.DeleteMessageBatch(ctx, &input)
			return err
		})

		return err
	}
}
