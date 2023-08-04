package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type RequestType string

const (
	subscribe   RequestType = "subscribe"
	unsubscribe RequestType = "unsubscribe"
)

type ChannelRequest struct {
	SubscriptionID string
	ChannelID      string
	RequestType    RequestType
	Params         *map[string]interface{}
}

func readWebsocket(socket *websocket.Conn, router *Router) {

	for {
		var request *ChannelRequest

		err := socket.ReadJSON(&request)
		if err != nil {
			return
		}

		channelID := request.ChannelID
		SubscriptionID := request.SubscriptionID

		switch request.RequestType {
		case subscribe:
			router.Subscribe(channelID, &Subscriber{
				socket:         socket,
				SubscriptionID: SubscriptionID,
			},
				request.Params,
			)
		case unsubscribe:
			router.Unsubscribe(SubscriptionID)
		}
	}
}

func HandleWs(writer http.ResponseWriter, request *http.Request, router *Router) {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(request *http.Request) bool {
			return true
		},
	}

	socket, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		return
	}

	defer func() {
		socket.Close()
		fmt.Println("closing")
	}()

	readWebsocket(socket, router)
}
