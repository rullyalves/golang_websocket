package messaging

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/messaging"
	"fmt"
)

type NotificationParams struct {
	Token    *string           `json:"token"`
	Topic    *string           `json:"topic"`
	Metadata map[string]string `json:"metadata"`
}

func (r NotificationParams) ToMessage() *messaging.Message {

	var topic string
	var token string

	if r.Topic != nil {
		topic = *r.Topic
	}

	if r.Token != nil {
		token = *r.Token
	}

	return &messaging.Message{
		Data:  r.Metadata,
		Topic: topic,
		Token: token,
	}
}

type SendNotification func(ctx context.Context, params []*NotificationParams) error

func GetFailedMessages(response *messaging.BatchResponse, messages []*messaging.Message) []*messaging.Message {
	var failedMessages []*messaging.Message

	if response.FailureCount > 0 {
		for idx, resp := range response.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedMessages = append(failedMessages, messages[idx])
			}
		}
	}

	return failedMessages
}

func SendPushNotification(client *messaging.Client) SendNotification {
	return func(ctx context.Context, params []*NotificationParams) error {

		var messages []*messaging.Message
		for _, value := range params {

			if value.Token == nil && value.Topic == nil {
				return errors.New("topic or Token must be not null")
			}
			messages = append(messages, value.ToMessage())
		}

		response, err := client.SendEach(ctx, messages)

		//TODO: Retry if failed message

		if response.FailureCount > 0 {
			failedMessages := GetFailedMessages(response, messages)
			fmt.Println(failedMessages)
		}

		return err
	}
}
