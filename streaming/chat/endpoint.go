package chat

import (
	"context"
	"fmt"
	"go_ws/shared/utils"
	chatDao "go_ws/streaming/chat/dao"
	"go_ws/streaming/chat/models"
	"go_ws/streaming/shared/broadcast"
	"go_ws/streaming/shared/websocket"

	"log"
)

func GetChatsAsStream(br broadcast.Observer[*models.ChatEvent], findChats chatDao.FindChatsEventsByUserId) websocket.EndpointHandler {
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
				data, err := findChats(context, *userId)
				if err != nil {
					log.Println(err)
				}

				if data != nil {
					//	sendData(data)
					fmt.Println(data)
				}
			}
		}()

		listener := func(message *models.ChatEvent) {
			if userId != nil && *userId == message.UserID {
				sendData(message)
			}
		}

		cancel := br.Subscribe(listener)
		return cancel
	}
}
