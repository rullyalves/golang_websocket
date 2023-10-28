package chat_service

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/chat_service/chat"
	chatDao "go_ws/services/chat_service/chat/dao/neo4j"
	"go_ws/services/chat_service/delivery"
	deliveryDao "go_ws/services/chat_service/delivery/dao/neo4j"
	deliveryUseCases "go_ws/services/chat_service/delivery/usecases"
	"go_ws/services/chat_service/message"
	messageDao "go_ws/services/chat_service/message/dao/neo4j"
	"go_ws/services/chat_service/message/usecases"
	"go_ws/shared/http_router"
)

func Handle(router http_router.Router, neo4jDriver *neo4j.DriverWithContext, sqsClient *sqs.Client) {
	// chats
	findAllChats := chatDao.FindByParticipantId(neo4jDriver)
	router.Get("/users/:userId/chats", chat.FindChatsHandler(findAllChats))

	// messages
	findMessagesById := messageDao.FindByMessageIdIn(neo4jDriver)
	findLastMessages := messageDao.FindLastMessages(neo4jDriver)
	findAllMessages := messageDao.FindAll(neo4jDriver)
	sendMessage := usecases.NewSendMessageCommand(neo4jDriver, sqsClient)

	router.Get("/messages", message.FindMessagesByIdInHandler(findMessagesById))
	router.Get("/messages/latest", message.FindLastMessagesHandler(findLastMessages))
	router.Get("/chats/:chatId/messages", message.FindMessagesHandler(findAllMessages))
	router.Post("/messages", message.SendTextMessageHandler(sendMessage))

	// deliveries
	findDeliveries := deliveryDao.FindByMessageIdIn(neo4jDriver)
	confirmDelivery := deliveryUseCases.NewConfirmMessageDelivery(neo4jDriver, sqsClient)

	router.Get("/deliveries", delivery.DeliveriesByMessageHandler(findDeliveries))
	router.Post("/deliveries", delivery.ConfirmMessageDeliveryStatusHandler(confirmDelivery))
}
