package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/chat_service/message/models"
	sharedModels "go_ws/services/chat_service/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type ChatLastMessagesResult struct {
	ChatID   string               `json:"chatId"`
	Messages []models.MessageView `json:"messages"`
}

type CreateMessageDataParams struct {
	ID        string                 `json:"id"`
	CreatedAt time.Time              `json:"createdAt"`
	Text      string                 `json:"text"`
	MediaType sharedModels.MediaType `json:"mediaType"`
	SenderId  string                 `json:"senderId"`
	ChatId    string                 `json:"chatId"`
	ParentId  *string                `json:"parentId"`
}

type FindLastMessagesByChatIdIn func(ctx context.Context, ids []string, limit int) (map[string][]models.MessageView, error)

type FindMessageById func(ctx context.Context, id string) (*models.MessageView, error)

type FindMessageByIdIn func(ctx context.Context, ids []string) (*[]models.MessageView, error)

type FindMessages func(ctx context.Context, chatId string, cursor time.Time, limit int) ([]models.MessageView, error)

type SaveMessage func(ctx context.Context, data CreateMessageDataParams) error

type DeleteMessage func(ctx context.Context, id string) error

func FindLastMessages(driver *neo4j.DriverWithContext) FindLastMessagesByChatIdIn {
	return func(ctx context.Context, ids []string, limit int) (map[string][]models.MessageView, error) {

		params := map[string]any{
			"ids":   ids,
			"limit": limit,
		}

		query := `
		MATCH (chat: Chat) WHERE chat.id IN ($ids)
		CALL {
			WITH chat
			MATCH (message: Message) WHERE (message)-->(chat) RETURN message ORDER BY message.createdAt DESC LIMIT $limit
		}
		WITH chat.id as ownerId, collect(message) as messages
		RETURN {chatId: ownerId, messages: messages}`

		result, err := neo4jdb.ExecuteWithMapping[ChatLastMessagesResult](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		mapResult := make(map[string][]models.MessageView)

		for _, item := range result {
			mapResult[item.ChatID] = item.Messages
		}

		return mapResult, nil
	}
}

func FindByMessageIdIn(driver *neo4j.DriverWithContext) FindMessageByIdIn {
	return func(ctx context.Context, ids []string) (*[]models.MessageView, error) {
		params := map[string]any{
			"ids": ids,
		}

		query := `
			MATCH (message: Message) WHERE message.id IN ($ids)
			RETURN message`

		results, err := neo4jdb.ExecuteWithMapping[models.MessageView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func FindAll(driver *neo4j.DriverWithContext) FindMessages {
	return func(ctx context.Context, chatId string, cursor time.Time, limit int) ([]models.MessageView, error) {

		params := map[string]any{
			"chatId": chatId,
			"limit":  limit,
			"cursor": cursor,
		}

		query := `
            MATCH (:Chat{id: $chatId})<--(m :Message)
			MATCH (sender: Profile)-->(m)
			WHERE m.createdAt < $cursor
			RETURN m ORDER BY m.createdAt DESC LIMIT $limit`

		results, err := neo4jdb.ExecuteWithMapping[models.MessageView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveMessage {
	return func(ctx context.Context, data CreateMessageDataParams) error {

		params := map[string]any{
			"id":        data.ID,
			"createdAt": data.CreatedAt,
			"text":      data.Text,
			"mediaType": data.MediaType,
			"senderId":  data.SenderId,
			"chatId":    data.ChatId,
			"parentId":  data.ParentId,
		}

		query := `
		MATCH (sender: Profile{id: $senderId})
		MATCH (chat: Chat{id: $chatId})
		WITH sender, chat
		MERGE (sender)-[:SEND]->(m: Message{id: $id})-[:TO]->(chat)
		SET
			m.text = $text,
			m.mediaType = $mediaType,
			m.createdAt = COALESCE($createdAt, m.createdAt),
			m.senderId = $senderId,
			m.chatId = $chatId,
			m.parentId = $parentId
		WITH m
		MATCH (parentMessage: Message{id: $parentId})
		MERGE (m)<-[:PARENT_OF]-(parentMessage)`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)

		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteMessage {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		_, err := neo4jdb.ExecuteQuery(ctx, driver, "MATCH (m :Message {id: $id}) DELETE m", params)
		return err
	}
}
