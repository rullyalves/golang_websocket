package websocket

import (
	"github.com/gorilla/websocket"
)

type Subscriber struct {
	SessionID string
	socket    *websocket.Conn
}

type Channel struct {
	subscribe   chan *Subscriber
	unsubscribe chan string
	broadcast   chan interface{}
	ChannelID   string
	subscribers map[string]*websocket.Conn
}

func NewChannel(channelID string) *Channel {
	return &Channel{
		ChannelID:   channelID,
		unsubscribe: make(chan string),
		subscribe:   make(chan *Subscriber),
		subscribers: map[string]*websocket.Conn{},
		broadcast:   make(chan interface{}),
	}
}

func (channel *Channel) Start(removeSession func(sessionID string)) {
	for {
		select {
		case message := <-channel.broadcast:
			for sessionID, socket := range channel.subscribers {
				if err := socket.WriteJSON(message); err != nil {
					removeSession(sessionID)
					socket.Close()
				}
			}
		case subscriber := <-channel.subscribe:
			socket := subscriber.socket
			channel.subscribers[subscriber.SessionID] = socket

			socket.SetCloseHandler(func(code int, text string) error {
				removeSession(subscriber.SessionID)
				return nil
			})

		case sessionID := <-channel.unsubscribe:
			delete(channel.subscribers, sessionID)
		}
	}
}

type Router struct {
	subscribe   chan Channel
	unsubscribe chan Channel
	channels    map[string]Channel
}

func NewRouter() *Router {
	return &Router{
		subscribe:   make(chan Channel),
		unsubscribe: make(chan Channel),
		channels:    map[string]Channel{},
	}
}

func (router *Router) Broadcast(channelID string, message interface{}) {
	router.channels[channelID].broadcast <- message
}

func (router *Router) Subscribe(channelID string, subscriber *Subscriber) {
	router.channels[channelID].subscribe <- subscriber
}

func (router *Router) Unsubscribe(channelID string, sessionID string) {
	router.channels[channelID].unsubscribe <- sessionID
}

func (router *Router) RemoveSubscriptions(sessionID string) {

	for _, channel := range router.channels {
		channel.unsubscribe <- sessionID
	}
}

func (router *Router) Register(channel Channel) {
	router.subscribe <- channel
}

func (router *Router) Start() {
	for {
		select {
		case channel := <-router.subscribe:
			router.channels[channel.ChannelID] = channel
			go channel.Start(router.RemoveSubscriptions)
		case channel := <-router.unsubscribe:
			delete(router.channels, channel.ChannelID)
		}
	}
}
