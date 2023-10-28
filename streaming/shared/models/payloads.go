package models

import (
	"github.com/google/uuid"
	shared "go_ws/shared/models"
	"go_ws/shared/queue"
	pushNotifications "go_ws/streaming/push_notifications"
	"maps"
	"time"
)

func CreateNotificationPayload(title string, body string, notificationType shared.NotificationType, topic string, extraData map[string]string) queue.MessageInput {

	metadata := map[string]string{
		"t":    title,
		"b":    body,
		"type": string(notificationType),
	}

	maps.Copy(metadata, extraData)

	data := pushNotifications.NotificationParams{
		Topic:    &topic,
		Metadata: metadata,
	}

	return queue.MessageInput{ID: uuid.NewString(), Body: data}
}

func CreateEventPayload[T any](data T, targetId string, createdAt time.Time, messageType shared.MessageType) queue.MessageInput {
	event := SubscriptionEvent[T]{
		ID:        uuid.NewString(),
		Data:      data,
		Type:      messageType,
		CreatedAt: createdAt,
		UserId:    targetId,
	}
	return queue.MessageInput{ID: uuid.NewString(), Body: event}
}
