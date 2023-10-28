package block

import (
	"context"
	blockDao "go_ws/streaming/block/dao"
	"go_ws/streaming/block/models"
	"go_ws/streaming/shared/broadcast"
	"go_ws/streaming/shared/websocket"
)

// TODO: carregar eventos não recebidos
func GetBlocksAsStream(br broadcast.Observer[*models.BlockEvent], findBlocks blockDao.FindBlockEventsByUserId) websocket.EndpointHandler {
	return func(context context.Context, sendData func(value interface{})) func() {

		listener := func(message *models.BlockEvent) {
			sendData(message)
		}

		cancel := br.Subscribe(listener)

		return cancel
	}
}
