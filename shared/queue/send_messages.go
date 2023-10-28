package queue

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go_ws/shared/utils"
	"log"
)

var errorWhenSendMessages = errors.New("error when send messages to queue")

type MessageInput struct {
	ID   string
	Body any
}

type SendToQueue func(context context.Context, queueUrl string, entries []MessageInput) error

type SendChunkToQueue[T any] func(ctx context.Context, messageQueueUrl string, slice []T, size int, mapTo MapToQueueMessage[T]) error

type MapToQueueMessage[T any] func(value T) (*MessageInput, error)

func SendMessageToQueue(client *sqs.Client) SendToQueue {
	return func(context context.Context, queueUrl string, inputs []MessageInput) error {

		var entries []types.SendMessageBatchRequestEntry

		for _, item := range inputs {
			id := item.ID

			jsonData, err := json.Marshal(item.Body)

			if err != nil {
				return err
			}

			body := string(jsonData)

			entries = append(entries, types.SendMessageBatchRequestEntry{
				Id:                     &id,
				MessageBody:            &body,
				MessageDeduplicationId: &id,
				MessageGroupId:         &id,
			})
		}

		data := sqs.SendMessageBatchInput{
			Entries:  entries,
			QueueUrl: &queueUrl,
		}
		result, err := client.SendMessageBatch(context, &data)

		if err != nil {
			return err
		}

		successful := result.Successful

		if len(entries) != len(successful) {
			return errorWhenSendMessages
		}

		return nil
	}
}

func SendMessagesToQueue[T any](sendToQueue SendToQueue) SendChunkToQueue[T] {
	return func(ctx context.Context, messageQueueUrl string, slice []T, size int, mapTo MapToQueueMessage[T]) error {
		chunks := utils.Chunk(slice, size)

		for _, chunk := range chunks {

			var entries []MessageInput

			for _, value := range chunk {

				entry, err := mapTo(value)

				if err != nil {
					return err
				}

				entries = append(entries, *entry)
			}

			err := sendToQueue(ctx, messageQueueUrl, entries)

			if err != nil {
				log.Print("error when send to queue")
				return err
			}
		}

		return nil
	}
}
