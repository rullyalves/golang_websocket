package like

import (
	"context"
	"go_ws/streaming/shared/broadcast"
	"go_ws/streaming/shared/websocket"

	likeDao "go_ws/streaming/like/dao"
	"go_ws/streaming/like/models"
)

// TODO: carregar eventos não recebidos
func GetLikesAsStream(br broadcast.Observer[*models.LikeEvent], findLikes likeDao.FindLikeEventsByUserId) websocket.EndpointHandler {
	return func(context context.Context, sendData func(value interface{})) func() {

		listener := func(message *models.LikeEvent) {

			sendData(message)
		}

		cancel := br.Subscribe(listener)

		return cancel
	}
}
