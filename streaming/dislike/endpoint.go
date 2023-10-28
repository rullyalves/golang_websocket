package dislike

import (
	"context"
	"go_ws/streaming/shared/broadcast"
	"go_ws/streaming/shared/websocket"

	dislikeDao "go_ws/streaming/dislike/dao"
	"go_ws/streaming/dislike/models"
)

// TODO: carregar eventos não recebidos
func GetDislikesAsStream(br broadcast.Observer[*models.DislikeEvent], findDislikes dislikeDao.FindDislikeEventsByUserId) websocket.EndpointHandler {
	return func(context context.Context, sendData func(value interface{})) func() {

		listener := func(message *models.DislikeEvent) {
			sendData(message)
		}

		cancel := br.Subscribe(listener)

		return cancel
	}
}
