package usecases

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	chatDao "go_ws/services/chat_service/chat/dao/neo4j"
	profileDao "go_ws/services/profile_service/profile/dao/neo4j"
	profileModels "go_ws/services/profile_service/profile/models"
	likeDao "go_ws/services/relationship_service/like/dao/neo4j"
	matchDao "go_ws/services/relationship_service/match/dao/neo4j"
	shared "go_ws/shared/models"
	"go_ws/shared/queue"
	streamingModels "go_ws/streaming/shared/models"
	"time"
)

type CreateLikeParams struct {
	ReceiverID string `json:"receiverId" validate:"omitempty,uuid4"`
}

type SendLike func(context.Context, string, CreateLikeParams) error

type ProfileWithName struct {
	ID   string
	Name string
}

func MapNames(items []profileModels.ProfileView) map[string]ProfileWithName {
	result := make(map[string]ProfileWithName)
	for _, item := range items {
		result[item.ID] = ProfileWithName{
			ID:   item.ID,
			Name: item.Name,
		}
	}

	return result
}

// TODO: enviar eventos de like para a fila de websocket e push
func sendLikeEvents(sendToQueue queue.SendToQueue) func(ctx context.Context, sender ProfileWithName, receiver ProfileWithName) error {
	return func(ctx context.Context, sender ProfileWithName, receiver ProfileWithName) error {

		var items []queue.MessageInput

		//TODO: mover para fila diferente
		likeNotification := streamingModels.CreateNotificationPayload(
			"Você recebeu um like",
			fmt.Sprintf("%v gostou de você !!", sender.Name),
			shared.NotificationTypeLike,
			fmt.Sprintf("users-%v", receiver.ID),
			map[string]string{
				"profile_id": sender.ID,
			},
		)

		items = append(items, likeNotification)

		likeEvent := streamingModels.CreateEventPayload[any](
			//TODO: pegar dados do like
			nil,
			receiver.ID,
			time.Now().UTC(),
			shared.MessageTypeUpdate,
		)

		items = append(items, likeEvent)

		//TODO: pegar url da fila do env
		return sendToQueue(ctx, queue.LikeStreamQueueUrl, items)
	}
}

func sendChatEvents(sendToQueue queue.SendToQueue) func(ctx context.Context, participantIds []string) error {
	return func(ctx context.Context, participantIds []string) error {
		var events []queue.MessageInput
		for _, participantId := range participantIds {

			chatEvent := streamingModels.CreateEventPayload[any](
				//TODO: pegar dados do chat
				nil,
				participantId,
				time.Now().UTC(),
				shared.MessageTypeUpdate,
			)

			events = append(events, chatEvent)
		}

		//TODO: por queue url correta do env

		return sendToQueue(ctx, queue.ChatStreamQueueUrl, events)
	}
}

func sendMatchEvents(sendToQueue queue.SendToQueue) func(ctx context.Context, sender ProfileWithName, receiver ProfileWithName) error {
	return func(ctx context.Context, sender ProfileWithName, receiver ProfileWithName) error {

		var items []queue.MessageInput

		notificationTitle := "Você deu Match"
		notificationBody := "Você deu Match com %v"
		topic := "users-%v"
		var notificationType shared.NotificationType = shared.NotificationTypeMatch

		senderNotification := streamingModels.CreateNotificationPayload(
			notificationTitle,
			fmt.Sprintf(notificationBody, receiver.Name),
			notificationType,
			fmt.Sprintf(topic, sender.ID),
			map[string]string{
				"profile_id": receiver.ID,
			},
		)

		items = append(items, senderNotification)

		receiverNotification := streamingModels.CreateNotificationPayload(
			notificationTitle,
			fmt.Sprintf(notificationBody, sender.Name),
			notificationType,
			fmt.Sprintf(topic, receiver.ID),
			map[string]string{
				"profile_id": sender.ID,
			},
		)

		items = append(items, receiverNotification)

		//TODO: por queue url correta do env

		return sendToQueue(ctx, queue.LikePushQueueUrl, items)
	}
}

func handleMatch(saveMatch matchDao.SaveMatch, saveChat chatDao.SaveChat, sendToQueue queue.SendToQueue) func(ctx context.Context, sender ProfileWithName, receiver ProfileWithName, likeIds []string) error {
	return func(ctx context.Context, sender ProfileWithName, receiver ProfileWithName, likeIds []string) error {
		participantIds := []string{receiver.ID, sender.ID}

		matchId := uuid.NewString()

		err := saveMatch(ctx, matchDao.CreateMatchDataParams{
			ID:        matchId,
			CreatedAt: time.Now().UTC(),
			LikeIds:   likeIds,
			IsActive:  true,
		})

		if err != nil {
			return err
		}

		err = saveChat(ctx, chatDao.CreateChatDataParams{
			ID:             matchId,
			CreatedAt:      time.Now().UTC().UTC(),
			MatchID:        matchId,
			ParticipantIDs: participantIds,
		})

		if err != nil {
			return err
		}

		//TODO: tomar cuidado se houver mais de 10 mensagens
		err = sendChatEvents(sendToQueue)(ctx, participantIds)
		if err != nil {
			return err
		}

		err = sendMatchEvents(sendToQueue)(ctx, sender, receiver)
		return err
	}
}

func sendLike(
	saveLike likeDao.SaveLike,
	saveMatch matchDao.SaveMatch,
	saveChat chatDao.SaveChat,
	findLikeBySenderAndReceiver likeDao.FindLikeBySenderAndReceiverId,
	findProfileByIdIn profileDao.FindProfileByIdIn,
	sendToQueue queue.SendToQueue,
) SendLike {

	return func(ctx context.Context, userId string, params CreateLikeParams) error {

		likeId := uuid.NewString()

		err := saveLike(ctx, userId, likeDao.CreateLikeDataParams{
			ID:         likeId,
			ReceiverID: params.ReceiverID,
			CreatedAt:  time.Now().UTC(),
		})

		if err != nil {
			return err
		}

		results, err := findProfileByIdIn(ctx, []string{userId, params.ReceiverID})

		if err != nil {
			return err
		}

		namesMap := MapNames(*results)

		sender := namesMap[userId]
		receiver := namesMap[params.ReceiverID]

		err = sendLikeEvents(sendToQueue)(ctx, sender, receiver)
		if err != nil {
			return err
		}

		inverseLike, err := findLikeBySenderAndReceiver(ctx, params.ReceiverID, userId)

		if inverseLike == nil {
			return nil
		}

		err = handleMatch(saveMatch, saveChat, sendToQueue)(ctx, sender, receiver, []string{inverseLike.ID, likeId})
		if err != nil {
			return err
		}

		return err
	}
}

func NewSendLike(driver *neo4j.DriverWithContext, sqsClient *sqs.Client) SendLike {
	saveLike := likeDao.Save(driver)
	saveMatch := matchDao.Save(driver)
	saveChat := chatDao.Save(driver)

	findLikeBySenderAndReceiver := likeDao.FindBySenderAndReceiverId(driver)
	findProfileByIdIn := profileDao.FindByIdIn(driver)
	sendToQueue := queue.SendMessageToQueue(sqsClient)

	return sendLike(saveLike, saveMatch, saveChat, findLikeBySenderAndReceiver, findProfileByIdIn, sendToQueue)
}
