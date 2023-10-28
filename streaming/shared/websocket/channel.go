package websocket

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	shared "go_ws/shared/models"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 3 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 3 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	subscribe   RequestType = "subscribe"
	unsubscribe RequestType = "unsubscribe"
)

type EndpointHandler func(context context.Context, sendData func(data interface{})) func()

type RequestType string

type SubscriptionResponse struct {
	SubscriptionID string             `json:"subscriptionId"`
	Operation      string             `json:"operation"`
	MessageType    shared.MessageType `json:"messageType"`
	Data           interface{}        `json:"data"`
}

type ChannelRequest struct {
	SubscriptionID string          `json:"subscriptionId"`
	Operation      string          `json:"operation"`
	RequestType    RequestType     `json:"requestType"`
	Params         map[string]*any `json:"params"`
}

type Subscriber struct {
	SubscriptionID string
	Operation      string
	socket         *websocket.Conn
}

type Router struct {
	handlers map[string]EndpointHandler
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]EndpointHandler),
	}
}

func (r *Router) Subscribe(handler EndpointHandler, subscriber *Subscriber, params map[string]*any, sendChannel chan interface{}) func() {

	newContext := context.Background()
	newContext = context.WithValue(newContext, "params", params)
	newContext = context.WithValue(newContext, "subscriptionId", subscriber.SubscriptionID)
	newContext = context.WithValue(newContext, "operation", subscriber.Operation)

	send := func(data interface{}) {

		sendChannel <- SubscriptionResponse{
			SubscriptionID: subscriber.SubscriptionID,
			Operation:      subscriber.Operation,
			MessageType:    shared.MessageTypeUpdate,
			Data:           data,
		}
	}

	dispose := handler(newContext, send)

	return dispose
}

func (r *Router) Handle(channelID string, handler EndpointHandler) {
	r.handlers[channelID] = handler
}

func (r *Router) handleMessages(socket *websocket.Conn, channel chan interface{}, ticker *time.Ticker) {
	for {
		select {
		case data, ok := <-channel:

			if !ok {
				return
			}

			err := socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("error when set write deadlines %v", err)
				return
			}

			err = socket.WriteJSON(data)
			if err != nil {
				log.Printf("error when: %v", err)
				return
			}

		case <-ticker.C:
			err := socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("error when set write deadline: %v", err)
			}

			if err := socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error when ping %v", err)
				return
			}
		}
	}
}

func (r *Router) handleDeadlines(ticker *time.Ticker, channel chan interface{}, socket *websocket.Conn) error {

	err := socket.SetReadDeadline(time.Now().Add(pongWait))

	if err != nil {
		log.Printf("error when set read deadline %v", err)
		return err
	}

	socket.SetPongHandler(func(string) error {
		err := socket.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Printf("error when set read deadline %v", err)
			return err
		}
		return nil
	})

	go r.handleMessages(socket, channel, ticker)

	return err
}

func (r *Router) unsubscribeAll(subscriptions map[string]func()) {
	for subscriptionID := range subscriptions {
		//log.Printf("remove subscription %v\n", key)
		dispose, ok := subscriptions[subscriptionID]

		if ok {
			dispose()
		}

		delete(subscriptions, subscriptionID)
	}
}

func (r *Router) unsubscribe(subscriptions map[string]func(), subscriptionID string) {
	dispose, ok := subscriptions[subscriptionID]

	if !ok {
		return
	}

	delete(subscriptions, subscriptionID)
	dispose()

	log.Printf("delete subscription: %v\n", subscriptions)
}

func (r *Router) startSubscription(subscriptions map[string]func(), request ChannelRequest, socket *websocket.Conn, channel chan interface{}) {

	operation := request.Operation
	subscriptionID := request.SubscriptionID

	handler, hasHandler := r.handlers[operation]

	if !hasHandler {
		fmt.Println("stop: not has handler")
		return
	}

	_, loaded := subscriptions[subscriptionID]

	if loaded {
		fmt.Println("stop: subscription already exists")
		return
	}

	subscriber := Subscriber{
		socket:         socket,
		SubscriptionID: subscriptionID,
		Operation:      operation,
	}

	dispose := r.Subscribe(
		handler,
		&subscriber,
		request.Params,
		channel,
	)

	subscriptions[subscriptionID] = dispose

	log.Printf("new subscription: %v\n", subscriptions)
}

func (r *Router) readWebsocket(socket *websocket.Conn) error {
	subscriptions := make(map[string]func())
	channel := make(chan interface{})
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		err := socket.Close()
		if err != nil {
			log.Printf("error when close socket %v", err)
		}
	}()

	err := r.handleDeadlines(ticker, channel, socket)
	if err != nil {
		return err
	}

	for {

		var request ChannelRequest

		err := socket.ReadJSON(&request)

		log.Printf("reading new data %v\n", request)

		if err != nil {
			r.unsubscribeAll(subscriptions)

			log.Printf("error while reading %v\n subscriptions: %v", err, subscriptions)

			return err
		}

		switch request.RequestType {
		case subscribe:
			r.startSubscription(subscriptions, request, socket, channel)
		case unsubscribe:
			r.unsubscribe(subscriptions, request.SubscriptionID)
		}
	}
}

func (r *Router) HandleWs(response http.ResponseWriter, request *http.Request) {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	socket, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Printf("upgrade error: %s", err)
		return
	}

	log.Println("receive a new connection")

	socketErr := r.readWebsocket(socket)

	if socketErr != nil {
		log.Printf("error when reading socket %v\n", socketErr)
		return
	}

	log.Println("finish connection")
}
