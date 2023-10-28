package main

import (
	"context"
	"fmt"
	"go_ws/services/authentication_service"
	"go_ws/services/chat_service"
	"go_ws/services/profile_service"
	"go_ws/services/relationship_service"
	"go_ws/services/shared/validations"
	"go_ws/shared/apis"
	"go_ws/shared/http_router"
	"go_ws/shared/queue"
	"go_ws/streaming"
	"go_ws/streaming/confirmation_service"
	pushNotifications "go_ws/streaming/push_notifications"
	"go_ws/streaming/shared/websocket"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8082"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	address := fmt.Sprintf("localhost:%v", port)

	ctx := context.Background()

	sqsClient, mongoDriver, neo4jDriver, err := apis.StartConnections(ctx)

	if err != nil {
		log.Fatalf("error while start database connections %v", err)
	}

	messagingClient, authClient, err := apis.StartFirebaseApis(ctx)

	if err != nil {
		log.Fatalf("error while start database connections %v", err)
	}

	validations.InitTranslations()

	wsRouter := websocket.NewRouter()

	router := http_router.New()

	fmt.Println(authClient)

	//	router.Use(authentication.AuthenticationMiddleware(authClient))

	router.Get("/ws", wsRouter.HandleWs)

	authentication_service.Handle(router, neo4jDriver, authClient)
	profile_service.Handle(router, neo4jDriver, sqsClient)
	chat_service.Handle(router, neo4jDriver, sqsClient)
	relationship_service.Handle(router, neo4jDriver, sqsClient)
	confirmation_service.Handle(router, mongoDriver)
	streaming.Handle(wsRouter, mongoDriver)

	processQueue := queue.ProcessQueue(sqsClient)
	startStreaming := streaming.StartStreamingConsumers(processQueue, mongoDriver)

	sendNotification := pushNotifications.SendPushNotification(messagingClient)
	sendMessage := pushNotifications.SendNotificationFromQueue(sendNotification)

	go pushNotifications.StartPushConsumers(processQueue, sendMessage)
	go startStreaming()

	log.Printf("connect to http://%s/", address)
	log.Fatal(http.ListenAndServe(address, router))

}
