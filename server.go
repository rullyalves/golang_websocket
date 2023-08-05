package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_ws/broadcast"
	"go_ws/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := gin.Default()

	wsRouter := websocket.NewRouter()

	ticker := time.NewTicker(500 * time.Millisecond)

	broadcaster := broadcast.NewBroadcaster()

	go func() {
		i := 0
		for {
			select {
			case <-ticker.C:
				broadcaster.Publish(i)
			}
			i++
		}
	}()

	wsRouter.Handle("deliveries", func(context context.Context) (chan interface{}, func()) {
		subscriptionId := context.Value("subscriptionId").(string)
		//params := context.Value("params").(*map[string]interface{})

		channel := make(chan interface{})

		listener := broadcast.Listener{
			ID: subscriptionId,
			OnData: func(message interface{}) {
				if !broadcast.IsCOpen(channel) {
					return
				}
				channel <- message
			},
		}
		broadcaster.Subscribe(listener)

		return channel, func() {
			broadcaster.Unsubscribe(listener)
		}
	})

	wsRouter.Handle("messages", func(context context.Context) (chan interface{}, func()) {
		subscriptionId := context.Value("subscriptionId").(string)
		//params := context.Value("params").(*map[string]interface{})
		//userId := (*params)["userId"]

		fmt.Println(subscriptionId)

		channel := make(chan interface{})

		return channel, nil
	})

	router.GET("/ws", func(context *gin.Context) {
		websocket.HandleWs(context.Writer, context.Request, wsRouter)
	})

	log.Printf("connect to http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
