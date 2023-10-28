package delivery

import (
	"context"
	"fmt"
	"go_ws/shared/utils"
	"go_ws/streaming/shared/broadcast"
	"go_ws/streaming/shared/websocket"

	deliveryDao "go_ws/streaming/delivery/dao"
	"go_ws/streaming/delivery/models"

	"log"
)

func GetDeliveriesAsStream(br broadcast.Observer[*models.DeliveryEvent], findDeliveries deliveryDao.FindDeliveriesByUserId) websocket.EndpointHandler {
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
				data, err := findDeliveries(context, *userId)
				if err != nil {
					log.Println(err)
				}

				if data != nil {
					//	sendData(data)
					fmt.Println(data)
				}
			}
		}()

		listener := func(message *models.DeliveryEvent) {
			if *userId == message.UserID {
				sendData(message)
			}
		}

		cancel := br.Subscribe(listener)

		return cancel
	}
}
