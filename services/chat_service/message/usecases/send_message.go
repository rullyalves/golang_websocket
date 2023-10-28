package usecases

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	messageDao "go_ws/services/chat_service/message/dao/neo4j"
	"go_ws/services/chat_service/message/models"
	"go_ws/services/chat_service/message/models/input"
	common "go_ws/services/chat_service/models"
	profileDao "go_ws/services/profile_service/profile/dao/neo4j"
	neo4jdb "go_ws/shared/database/neo4j"
	shared "go_ws/shared/models"
	"go_ws/shared/queue"
	streamingModels "go_ws/streaming/shared/models"
	"time"
)

type SendMessage func(ctx context.Context, body input.CreateTextMessageParamsInput) error

func NewSendMessageCommand(neo4jDriver *neo4j.DriverWithContext, sqsClient *sqs.Client) SendMessage {
	execute := sendMessageCommand(
		queue.SendMessageToQueue(sqsClient),
		messageDao.Save(neo4jDriver),
		profileDao.FindIdsByChatIdIn(neo4jDriver),
		profileDao.FindNameById(neo4jDriver),
	)
	txManager := neo4jdb.WithTransaction[any](neo4jDriver)

	return func(ctx context.Context, body input.CreateTextMessageParamsInput) error {
		_, err := txManager(ctx, func(ctx context.Context) (*any, error) {
			return nil, execute(ctx, body)
		})

		return err
	}
}

func sendMessageCommand(
	sendToQueue queue.SendToQueue,
	saveMessage messageDao.SaveMessage,
	findParticipants profileDao.FindParticipantIdsByChatIdIn,
	findSenderName profileDao.FindProfileNameById,
) SendMessage {

	return func(ctx context.Context, body input.CreateTextMessageParamsInput) error {

		message := body.Message
		chatID := message.ChatID
		senderID := message.SenderID

		newData := messageDao.CreateMessageDataParams{
			ID:        message.ID,
			CreatedAt: message.CreatedAt,
			Text:      body.Text,
			MediaType: common.MediaTypeText,
			ChatId:    chatID,
			ParentId:  message.ParentMessageID,
			SenderId:  message.SenderID,
		}

		err := saveMessage(ctx, newData)

		if err != nil {
			return err
		}

		senderName, err := findSenderName(ctx, senderID)

		if err != nil {
			return err
		}

		participantIds, err := findParticipants(ctx, []string{chatID})

		if err != nil {
			return err
		}

		err = sendEvents(sendToQueue)(ctx, chatID, *participantIds, *senderName, models.MessageView{
			ID:        message.ID,
			CreatedAt: message.CreatedAt,
			Text:      &body.Text,
			MediaType: common.MediaTypeText,
			ChatID:    chatID,
			ParentID:  message.ParentMessageID,
			SenderID:  senderID,
		})

		if err != nil {
			return err
		}

		return sendNotifications(sendToQueue)(ctx, *participantIds, senderID, *senderName, body.Text)
	}
}

func sendEvents(sendToQueue queue.SendToQueue) func(ctx context.Context, chatId string, participantIds []string, senderName string, message models.MessageView) error {

	return func(ctx context.Context, chatId string, participantIds []string, senderName string, message models.MessageView) error {

		var events []queue.MessageInput
		for _, participantId := range participantIds {

			messageEvent := streamingModels.CreateEventPayload[any](
				message,
				participantId,
				time.Now().UTC(),
				shared.MessageTypeUpdate,
			)

			events = append(events, messageEvent)
		}

		return sendToQueue(ctx, queue.MessageStreamQueueUrl, events)
	}
}

func sendNotifications(sendToQueue queue.SendToQueue) func(ctx context.Context, participantIds []string, senderId string, senderName string, text string) error {

	return func(ctx context.Context, participantIds []string, senderId string, senderName string, text string) error {

		var events []queue.MessageInput
		for _, participantId := range participantIds {

			if participantId == senderId {
				continue
			}

			messageEvent := streamingModels.CreateNotificationPayload(
				senderName,
				text,
				shared.NotificationTypeMessage,
				fmt.Sprintf("users-%v", participantId),
				map[string]string{
					"profile_id": senderId,
				},
			)

			events = append(events, messageEvent)
		}

		return sendToQueue(ctx, queue.MessagePushQueueUrl, events)
	}
}
