package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_ws/websocket"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const defaultPort = "8080"

type Broadcaster struct {
	listeners sync.Map
}

func isOpen(ch <-chan interface{}) bool {
	select {
	case <-ch:
		return false
	default:
		return true
	}
}

func (b *Broadcaster) Publish(message interface{}) {

	b.listeners.Range(func(key, _ any) bool {
		channel := key.(chan interface{})

		if !isOpen(channel) {
			b.Unsubscribe(channel)
			return true
		}

		channel <- message
		return true
	})
}

func (b *Broadcaster) Subscribe(channel chan interface{}) {
	b.listeners.Store(channel, nil)
}

func (b *Broadcaster) Unsubscribe(channel chan interface{}) {
	b.listeners.Delete(channel)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := gin.Default()

	wsRouter := websocket.NewRouter()

	ticker := time.NewTicker(500 * time.Millisecond)

	broadcast := Broadcaster{
		listeners: sync.Map{},
	}

	go func() {
		i := 0
		for {
			select {
			case <-ticker.C:
				fmt.Println("ticker")
				broadcast.Publish(i)
			}
			i++
		}
	}()

	wsRouter.Handle("deliveries", func(context context.Context) chan interface{} {
		channel := make(chan interface{})
		broadcast.Subscribe(channel)
		return channel
	})

	wsRouter.Handle("messages", func(context context.Context) chan interface{} {
		channel := make(chan interface{})
		broadcast.Subscribe(channel)
		return channel
	})

	router.GET("/ws", func(context *gin.Context) {
		websocket.HandleWs(context.Writer, context.Request, wsRouter)
	})

	log.Printf("connect to http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
