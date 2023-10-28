package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/profile_service/preferences/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreatePreferencesDataParams struct {
	ID        string     `json:"id"`
	MinAge    *int       `json:"min_age"`
	MaxAge    *int       `json:"max_age"`
	Distance  *int       `json:"distance"`
	CreatedAt *time.Time `json:"createdAt"`
}

type FindPreferencesByProfileId func(context.Context, string) (*models.PreferencesView, error)

type FindPreferencesByProfileIdIn func(context.Context, []string) ([]models.PreferencesOwnerView, error)

type SavePreferences func(context.Context, string, CreatePreferencesDataParams) error

func FindByProfileId(driver *neo4j.DriverWithContext) FindPreferencesByProfileId {
	return func(ctx context.Context, id string) (*models.PreferencesView, error) {

		result, err := FindByProfileIdIn(driver)(ctx, []string{id})

		if err != nil {
			return nil, err
		}

		if len(result) == 0 {
			return nil, nil
		}

		return &result[0].Preferences, nil
	}
}

func FindByProfileIdIn(driver *neo4j.DriverWithContext) FindPreferencesByProfileIdIn {
	return func(ctx context.Context, ids []string) ([]models.PreferencesOwnerView, error) {
		params := map[string]any{
			"ids": ids,
		}

		query :=
			"MATCH (p: Profile) WHERE p.id IN ($ids) " +
				"MATCH (p)-->(pref :Preferences) " +
				"RETURN {ownerId: p.id, data: pref}"

		results, err := neo4jdb.ExecuteWithMapping[models.PreferencesOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SavePreferences {
	return func(ctx context.Context, userId string, data CreatePreferencesDataParams) error {
		params := map[string]any{
			"id":        data.ID,
			"minAge":    data.MinAge,
			"maxAge":    data.MaxAge,
			"distance":  data.Distance,
			"createdAt": data.CreatedAt,
			"profileId": userId,
		}

		query := `
		MATCH (p: Profile{id: $profileId}) 
		MERGE (p)-[:HAS]->(pref :Preferences) 
		ON CREATE SET
			pref.id = COALESCE($id, pref.id),
			pref.createdAt = COALESCE($createdAt, pref.createdAt)
		SET
			pref.minAge = COALESCE($minAge, pref.minAge),
			pref.maxAge = COALESCE($maxAge, pref.maxAge),
			pref.distance = COALESCE($distance, pref.distance)
		RETURN pref.id`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)

		return err
	}
}
