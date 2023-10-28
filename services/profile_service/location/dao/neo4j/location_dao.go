package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/profile_service/location/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateLocationDataParams struct {
	ID        *string    `json:"id"`
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	CreatedAt *time.Time `json:"createdAt"`
	ProfileID string     `json:"profileId"`
}

type UpdateLocationDataParams struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type FindLocationByProfileId func(ctx context.Context, id string) (*models.LocationView, error)

type FindLocationByProfileIdIn func(ctx context.Context, ids []string) (*[]models.LocationOwnerView, error)

type SaveLocation func(context.Context, CreateLocationDataParams) error

type DeleteLocationById func(ctx context.Context, id string) error

func FindProfileId(driver *neo4j.DriverWithContext) FindLocationByProfileId {
	return func(ctx context.Context, id string) (*models.LocationView, error) {

		result, err := FindProfileIdIn(driver)(ctx, []string{id})

		if err != nil {
			return nil, err
		}

		data := *result

		if len(data) == 0 {
			return nil, err
		}

		return &data[0].Location, err
	}
}

func FindProfileIdIn(driver *neo4j.DriverWithContext) FindLocationByProfileIdIn {
	return func(ctx context.Context, ids []string) (*[]models.LocationOwnerView, error) {
		params := map[string]any{
			"ids": ids,
		}

		query := `
		MATCH (p:Profile) WHERE p.id IN ($ids) 
		MATCH (p)-->(a:Location)
		RETURN {ownerId: p.id, Location: a}`

		results, err := neo4jdb.ExecuteWithMapping[models.LocationOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveLocation {
	return func(ctx context.Context, data CreateLocationDataParams) error {

		params := map[string]any{
			"profileId": data.ProfileID,
			"id":        data.ID,
			"latitude":  data.Latitude,
			"longitude": data.Longitude,
		}

		createdAt := data.CreatedAt

		if createdAt != nil {
			params["createdAt"] = *createdAt
		}

		query := `
		MATCH (p:Profile {id: $profileId}) 
		MERGE (p)-[:IS_AT]->(a:Location) 
		ON CREATE SET
			a.id = $id,
			a.createdAt = COALESCE($createdAt, a.createdAt)
		SET
			a.latitude = $latitude,
			a.longitude = $longitude`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteLocationById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		_, err := neo4jdb.ExecuteQuery(ctx, driver, "MATCH (a:Location {id: $id}) DELETE a", params)
		return err
	}
}
