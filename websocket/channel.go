package websocket

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type Subscriber struct {
	SubscriptionID string
	socket         *websocket.Conn
}

type Router struct {
	handlers      sync.Map
	subscriptions sync.Map
}

func NewRouter() *Router {
	return &Router{
		handlers:      sync.Map{},
		subscriptions: sync.Map{},
	}
}

func (router *Router) Unsubscribe(subscriptionID string) {
	channel, ok := router.subscriptions.Load(subscriptionID)
	if !ok {
		return
	}
	router.subscriptions.Delete(subscriptionID)
	close(channel.(chan interface{}))
}

func (router *Router) Subscribe(channelID string, subscriber *Subscriber) {

	handler, ok := router.handlers.Load(channelID)

	if !ok {
		return
	}

	channel := handler.(func(context.Context) chan interface{})(context.Background())

	router.subscriptions.Store(subscriber.SubscriptionID, channel)

	go router.Listen(channel, subscriber)
}

func (router *Router) Listen(channel chan interface{}, subscriber *Subscriber) {
	socket := subscriber.socket

	socket.SetCloseHandler(func(code int, text string) error {
		fmt.Println("closing socket")
		close(channel)
		return nil
	})

	for {
		select {
		case value, ok := <-channel:
			if !ok {
				fmt.Println("closed channel - request closing socket")
				return
			}
			err := socket.WriteJSON(value)
			if err != nil {
				fmt.Println(err)
				fmt.Println("request closing socket")
				socket.Close()
				return
			}
		}
	}
}

func (router *Router) Handle(channelID string, handler func(context context.Context) chan interface{}) {
	router.handlers.Store(channelID, handler)
}
