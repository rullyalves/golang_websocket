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
	broadcast   RequestType = "broadcast"
)

type ChannelRequest struct {
	SessionID   string
	ChannelID   string
	RequestType RequestType
}

func readWebsocket(socket *websocket.Conn, router *Router) {

	for {
		var request *ChannelRequest

		err := socket.ReadJSON(&request)
		if err != nil {
			return
		}

		channelID := request.ChannelID
		sessionID := request.SessionID

		switch request.RequestType {
		case subscribe:
			var filter Filter = func(message interface{}) bool {

				fmt.Println(message)

				return true
			}
			router.Subscribe(channelID, &Subscriber{
				socket:    socket,
				SessionID: sessionID,
				filter:    &filter,
			})
		case unsubscribe:
			router.Unsubscribe(channelID, sessionID)
		case broadcast:
			router.Broadcast(channelID, request)
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
