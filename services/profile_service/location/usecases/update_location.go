package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	locationDao "go_ws/services/profile_service/location/dao/neo4j"
	"time"
)

type CreateLocationByLocationParams struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type UpdateProfileLocation func(ctx context.Context, profileId string, params CreateLocationByLocationParams) error

func NewUpdateLocation(neo4jDriver *neo4j.DriverWithContext) UpdateProfileLocation {
	saveLocation := locationDao.Save(neo4jDriver)
	find := locationDao.FindProfileId(neo4jDriver)
	updateLocation := UpdateLocation(saveLocation, find)
	return updateLocation
}

func UpdateLocation(
	saveLocation locationDao.SaveLocation,
	find locationDao.FindLocationByProfileId,
) UpdateProfileLocation {

	return func(ctx context.Context, profileId string, params CreateLocationByLocationParams) error {

		currentLocation, err := find(ctx, profileId)

		if err != nil {
			return err
		}

		var LocationId string
		var createdAt time.Time

		if currentLocation != nil {
			LocationId = currentLocation.ID
			createdAt = currentLocation.CreatedAt
		} else {
			LocationId = uuid.NewString()
			createdAt = time.Now().UTC()
		}

		newLocation := locationDao.CreateLocationDataParams{
			ID:        &LocationId,
			Longitude: params.Longitude,
			Latitude:  params.Latitude,
			CreatedAt: &createdAt,
			ProfileID: profileId,
		}

		return saveLocation(ctx, newLocation)
	}
}
