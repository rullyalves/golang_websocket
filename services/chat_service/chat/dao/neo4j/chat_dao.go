package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/chat_service/chat/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateChatDataParams struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	MatchID        string    `json:"matchId"`
	ParticipantIDs []string  `json:"participantIds"`
}

type FindChatsByParticipantId func(ctx context.Context, id string, limit int, cursor time.Time) ([]models.ChatView, error)

type FindChats func(ctx context.Context) ([]models.ChatView, error)

type FindChatById func(ctx context.Context, id string) (*models.ChatView, error)

type SaveChat func(ctx context.Context, data CreateChatDataParams) error

type DeleteChat func(ctx context.Context, id string) error

func FindByParticipantId(driver *neo4j.DriverWithContext) FindChatsByParticipantId {
	return func(ctx context.Context, id string, limit int, cursor time.Time) ([]models.ChatView, error) {
		params := map[string]any{
			"id":     id,
			"limit":  limit,
			"cursor": cursor,
		}

		query := `
		MATCH (Profile{id: $id})-->(chat:Chat) 
		WHERE chat.createdAt < $cursor 
		RETURN chat 
		ORDER BY chat.createdAt DESC 
		LIMIT $limit`

		results, err := neo4jdb.ExecuteWithMapping[models.ChatView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func FindById(driver *neo4j.DriverWithContext) FindChatById {
	return func(ctx context.Context, id string) (*models.ChatView, error) {
		params := map[string]any{
			"id": id,
		}

		query := "MATCH (ch :Chat {id: $id}) RETURN ch"

		result, err := neo4jdb.ExecuteWithMapping[models.ChatView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		if len(result) == 0 {
			return nil, nil
		}

		return &result[0], nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveChat {
	return func(ctx context.Context, data CreateChatDataParams) error {
		params := map[string]any{
			"id":             data.ID,
			"createdAt":      data.CreatedAt,
			"matchId":        data.MatchID,
			"participantIds": data.ParticipantIDs,
		}

		query := `
		MATCH (participant :Profile) WHERE participant.id IN ($participantIds) 
		MATCH (mh: Match{id: $matchId}) 
		WITH participant, mh 
		MERGE (mh)-[:CREATE]->(ch: Chat) 
		ON CREATE SET 
			ch.id = $id,
			ch.matchId = $matchId,
			ch.participantIds = $participantIds,
			ch.createdAt = COALESCE($createdAt, ch.createdAt) 
		MERGE (participant)-[:PARTICIPATES]->(ch)`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteChat {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		query := "MATCH (ch :Chat {id: $id}) DELETE ch"
		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}
