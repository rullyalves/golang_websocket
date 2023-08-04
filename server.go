package main

import (
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

	channelRouter := websocket.NewRouter()

	go channelRouter.Start()

	channelRouter.Register(
		*websocket.NewChannel("2"),
	)

	router.GET("/ws", func(context *gin.Context) {
		websocket.HandleWs(context.Writer, context.Request, channelRouter)
	})

	log.Printf("connect to http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
