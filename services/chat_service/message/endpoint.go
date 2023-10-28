package message

import (
	messageDao "go_ws/services/chat_service/message/dao/neo4j"
	"go_ws/services/chat_service/message/models"
	"go_ws/services/chat_service/message/models/input"
	"go_ws/services/chat_service/message/usecases"
	"go_ws/shared/http_router"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func FindMessagesByIdInHandler(findMessages messageDao.FindMessageByIdIn) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[[]models.MessageView]

		execute(writer, func() (*[]models.MessageView, error) {
			vars := http_router.Vars(request)
			ids := strings.Split(vars.Get("ids"), ",")

			result, err := findMessages(request.Context(), ids)

			return result, err
		})
	}
}

func FindLastMessagesHandler(findMessages messageDao.FindLastMessagesByChatIdIn) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[map[string][]models.MessageView]

		execute(writer, func() (*map[string][]models.MessageView, error) {
			vars := http_router.Vars(request)
			chatIds := strings.Split(vars.Get("chatIds"), ",")
			limit, _ := strconv.Atoi(vars.Get("limit"))

			result, err := findMessages(request.Context(), chatIds, limit)

			return &result, err
		})
	}
}

func FindMessagesHandler(findMessages messageDao.FindMessages) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[[]models.MessageView]

		execute(writer, func() (*[]models.MessageView, error) {

			vars := http_router.Vars(request)
			params := http_router.Params(request)

			chatId := params["chatId"].(string)
			limit, _ := strconv.Atoi(vars.Get("limit"))
			cursor, _ := time.Parse(time.RFC3339, vars.Get("cursor"))

			result, err := findMessages(request.Context(), chatId, cursor, limit)

			return &result, err
		})
	}
}

func SendTextMessageHandler(sendMessage usecases.SendMessage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteMutationRequest[any, input.CreateTextMessageParamsInput]

		execute(writer, request, func(body input.CreateTextMessageParamsInput) (*any, error) {
			return nil, sendMessage(request.Context(), body)
		})

		writer.WriteHeader(http.StatusCreated)
	}
}
