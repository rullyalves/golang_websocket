package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_ws/websocket"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := gin.Default()

	wsRouter := websocket.NewRouter()

	wsRouter.Handle("deliveries", func(context context.Context) chan interface{} {
		channel := make(chan interface{})
		go func() {
			channel <- "hello man"
		}()
		return channel
	})

	wsRouter.Handle("messages", func(context context.Context) chan interface{} {
		channel := make(chan interface{})
		go func() {
			channel <- "hello man"
		}()
		return channel
	})

	router.GET("/ws", func(context *gin.Context) {
		websocket.HandleWs(context.Writer, context.Request, wsRouter)
	})

	log.Printf("connect to http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
