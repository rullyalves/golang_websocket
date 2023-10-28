package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/profile_service/address/models"
	neo4jdb "go_ws/shared/database/neo4j"
	"time"
)

type CreateAddressDataParams struct {
	ID           *string    `json:"id"`
	City         string     `json:"city"`
	Neighborhood string     `json:"neighborhood"`
	State        string     `json:"state"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	CreatedAt    *time.Time `json:"createdAt"`
	ProfileID    string     `json:"profileId"`
}

type UpdateAddressDataParams struct {
	ID           string    `json:"id"`
	City         string    `json:"city"`
	Neighborhood string    `json:"neighborhood"`
	State        string    `json:"state"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	CreatedAt    time.Time `json:"createdAt"`
}

type FindAddressByProfileId func(ctx context.Context, id string) (*models.AddressView, error)

type FindAddressByProfileIdIn func(ctx context.Context, ids []string) (*[]models.AddressOwnerView, error)

type SaveAddress func(context.Context, CreateAddressDataParams) error

type DeleteAddressById func(ctx context.Context, id string) error

func FindProfileId(driver *neo4j.DriverWithContext) FindAddressByProfileId {
	return func(ctx context.Context, id string) (*models.AddressView, error) {

		result, err := FindProfileIdIn(driver)(ctx, []string{id})

		if err != nil {
			return nil, err
		}

		data := *result

		if len(data) == 0 {
			return nil, err
		}

		return &data[0].Address, err
	}
}

func FindProfileIdIn(driver *neo4j.DriverWithContext) FindAddressByProfileIdIn {
	return func(ctx context.Context, ids []string) (*[]models.AddressOwnerView, error) {
		params := map[string]any{
			"ids": ids,
		}

		query := `
		MATCH (p:Profile) WHERE p.id IN ($ids) 
		MATCH (p)-->(a:Address)
		RETURN {ownerId: p.id, address: a}`

		results, err := neo4jdb.ExecuteWithMapping[models.AddressOwnerView](ctx, driver, query, params)
		if err != nil {
			return nil, err
		}

		return &results, nil
	}
}

func Save(driver *neo4j.DriverWithContext) SaveAddress {
	return func(ctx context.Context, data CreateAddressDataParams) error {

		params := map[string]any{
			"profileId":    data.ProfileID,
			"id":           data.ID,
			"city":         data.City,
			"state":        data.State,
			"neighborhood": data.Neighborhood,
			"latitude":     data.Latitude,
			"longitude":    data.Longitude,
		}

		createdAt := data.CreatedAt

		if createdAt != nil {
			params["createdAt"] = *createdAt
		}

		query := `
		MATCH (p:Profile {id: $profileId}) 
		MERGE (p)-[:IS_AT]->(a:Address) 
		ON CREATE SET
			a.id = $id,
			a.createdAt = COALESCE($createdAt, a.createdAt)
		SET 
			a.city = $city,
			a.neighborhood = $neighborhood,
			a.state = $state,
			a.latitude = $latitude,
			a.longitude = $longitude`

		_, err := neo4jdb.ExecuteQuery(ctx, driver, query, params)
		return err
	}
}

func Delete(driver *neo4j.DriverWithContext) DeleteAddressById {
	return func(ctx context.Context, id string) error {
		params := map[string]any{"id": id}
		_, err := neo4jdb.ExecuteQuery(ctx, driver, "MATCH (a:Address {id: $id}) DELETE a", params)
		return err
	}
}
