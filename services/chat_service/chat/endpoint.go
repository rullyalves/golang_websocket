package chat

import (
	chatDao "go_ws/services/chat_service/chat/dao/neo4j"
	"go_ws/services/chat_service/chat/models"
	"go_ws/shared/http_router"
	"net/http"
	"strconv"
	"time"
)

func FindChatsHandler(findChats chatDao.FindChatsByParticipantId) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		execute := http_router.ExecuteReadRequest[[]models.ChatView]

		execute(writer, func() (*[]models.ChatView, error) {

			vars := request.URL.Query()
			params := http_router.Params(request)

			participantId := params["userId"].(string)

			cursor, err := time.Parse(time.RFC3339, vars.Get("cursor"))

			if err != nil {
				return nil, err
			}

			limit, err := strconv.Atoi(vars.Get("limit"))

			if err != nil {
				return nil, err
			}

			chats, err := findChats(request.Context(), participantId, limit, cursor)

			return &chats, err
		})
	}
}
