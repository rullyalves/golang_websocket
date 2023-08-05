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

	channel, ok := router.subscriptions.LoadAndDelete(subscriptionID)
	if !ok {
		return
	}
	close(channel.(chan interface{}))
}

func (router *Router) Subscribe(channelID string, subscriber *Subscriber, params *map[string]interface{}) {

	handler, ok := router.handlers.Load(channelID)

	if !ok {
		return
	}

	_, loaded := router.subscriptions.Load(subscriber.SubscriptionID)
	if loaded {
		return
	}

	newContext := context.Background()
	newContext = context.WithValue(newContext, "params", params)
	newContext = context.WithValue(newContext, "subscriptionId", subscriber.SubscriptionID)

	channel, dispose := handler.(func(context context.Context) (chan interface{}, func()))(newContext)

	router.subscriptions.Store(subscriber.SubscriptionID, channel)

	go router.Listen(channel, subscriber, dispose)
}

func (router *Router) Listen(channel chan interface{}, subscriber *Subscriber, dispose func()) {

	defer func() {
		router.subscriptions.Delete(subscriber.SubscriptionID)

		if dispose != nil {
			dispose()
		}

		fmt.Println("defer listen")
	}()

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
				fmt.Println("closed channel")
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

func (router *Router) Handle(channelID string, handler func(context context.Context) (chan interface{}, func())) {
	router.handlers.Store(channelID, handler)
}
