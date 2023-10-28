package message

import (
	"context"
	"fmt"
	"go_ws/shared/utils"
	messageDao "go_ws/streaming/message/dao"
	"go_ws/streaming/message/models"
	"go_ws/streaming/shared/broadcast"
	"go_ws/streaming/shared/websocket"
	"log"
)

// TODO: carregar eventos não recebidos
func GetMessagesAsStream(br broadcast.Observer[*models.MessageEvent], findMessages messageDao.FindMessageEventsByUserId) websocket.EndpointHandler {
	return func(context context.Context, sendData func(value interface{})) func() {

		params := context.Value("params").(map[string]*any)
		subscriptionId := context.Value("subscriptionId").(string)

		var userId *string

		if params != nil {
			userId = utils.GetFromMap[string](params, "userId")
		}

		fmt.Println(subscriptionId)

		go func() {
			if userId != nil {
				data, err := findMessages(context, *userId)
				if err != nil {
					log.Println(err)
				}

				if data != nil {
					fmt.Println(data)
				}
			}
		}()

		listener := func(message *models.MessageEvent) {
			if *userId == message.UserId {
				sendData(message)
			}
		}

		cancel := br.Subscribe(listener)

		return cancel
	}
}
