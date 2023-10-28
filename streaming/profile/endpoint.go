package profile

import (
	"context"
	profileDao "go_ws/streaming/profile/dao"
	"go_ws/streaming/profile/models"
	"go_ws/streaming/shared/broadcast"
	"go_ws/streaming/shared/websocket"
)

// TODO: carregar eventos não recebidos
func GetProfilesAsStream(br broadcast.Observer[*models.ProfileEvent], findProfiles profileDao.FindProfileEventsByUserId) websocket.EndpointHandler {
	return func(context context.Context, sendData func(value interface{})) func() {

		listener := func(message *models.ProfileEvent) {
			sendData(message)
		}

		cancel := br.Subscribe(listener)

		return cancel
	}
}
