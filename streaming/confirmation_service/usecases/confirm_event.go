package usecases

import (
	"context"
	eventDao "go_ws/streaming/confirmation_service/dao/mongodb"
)

type ConfirmSubscriptionEvent func(context.Context, string, eventDao.Event) error

func ConfirmEvent(deleteEvent eventDao.DeleteEventById) ConfirmSubscriptionEvent {
	return func(ctx context.Context, eventId string, event eventDao.Event) error {
		return deleteEvent(ctx, eventId, event)
	}
}
