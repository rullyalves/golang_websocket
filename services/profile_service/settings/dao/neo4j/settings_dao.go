package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/profile_service/settings/models"

	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateSettingsDataParams struct {
	ID                             string    `json:"id"`
	IsVisible                      *bool     `json:"isVisible"`
	AllowReceiveMatchNotifications *bool     `json:"allowReceiveMatchNotifications"`
	AllowReceiveLikeNotifications  *bool     `json:"allowReceiveLikeNotifications"`
	CreatedAt                      time.Time `json:"created_at"`
}

type FindSettingsByProfileIdIn func(context.Context, []string) ([]models.SettingsOwnerView, error)

type SaveSettings func(context.Context, string, CreateSettingsDataParams) error

func FindProfileIdIn(driver *neo4j.DriverWithContext) FindSettingsByProfileIdIn {
	return func(ctx context.Context, ids []string) ([]models.SettingsOwnerView, error) {
		params := map[string]any{
			"ids": ids,
		}

		query := "MATCH (p:Profile)-[h]->(s:Settings) WHERE p.id IN ($ids) RETURN {ownerId: p.id, data: s}"

		results, err := neo4jdb.ExecuteWithMapping[models.SettingsOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveSettings {
	return func(ctx context.Context, userId string, data CreateSettingsDataParams) error {

		params := map[string]any{
			"id":                             data.ID,
			"isVisible":                      data.IsVisible,
			"allowReceiveMatchNotifications": data.AllowReceiveMatchNotifications,
			"allowReceiveLikeNotifications":  data.AllowReceiveLikeNotifications,
			"createdAt":                      data.CreatedAt,
			"profileId":                      userId,
		}

		//TODO: verificar se essa query possibilita mandar mais de um feedback por usuário

		query := `
		MATCH (profile: Profile{id: $profileId})
		MERGE (profile)-[:HAS]->(s: Settings)
		ON CREATE SET
			s.id = $id,
			s.createdAt = COALESCE($createdAt, s.createdAt) 
		SET 
			s.isVisible = COALESCE($isVisible, s.isVisible), 
			s.allowReceiveMatchNotifications = COALESCE($allowReceiveMatchNotifications, s.allowReceiveMatchNotifications),
			s.allowReceiveLikeNotifications = COALESCE($allowReceiveLikeNotifications, s.allowReceiveLikeNotifications)`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}
