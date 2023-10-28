package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/profile_service/address/api"
	addressDao "go_ws/services/profile_service/address/dao/neo4j"
	"time"
)

type CreateAddressByLocationParams struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type UpdateProfileLocation func(ctx context.Context, profileId string, params CreateAddressByLocationParams) error

func NewUpdateLocation(neo4jDriver *neo4j.DriverWithContext) UpdateProfileLocation {
	saveAddress := addressDao.Save(neo4jDriver)
	find := addressDao.FindProfileId(neo4jDriver)
	updateLocation := UpdateLocation(saveAddress, find, api.FindAddress)
	return updateLocation
}

func UpdateLocation(
	saveAddress addressDao.SaveAddress,
	find addressDao.FindAddressByProfileId,
	findAddressGeocoding api.FindAddressByLatLng,
) UpdateProfileLocation {

	return func(ctx context.Context, profileId string, params CreateAddressByLocationParams) error {

		//TODO: adicionar autenticação e autorização
		//TODO: verificar se a latitude é float64 mesmo

		apiResult, err := findAddressGeocoding(params.Latitude, params.Longitude)

		if err != nil {
			return err
		}

		currentAddress, err := find(ctx, profileId)

		if err != nil {
			return err
		}

		var addressId string
		var createdAt time.Time

		if currentAddress != nil {
			addressId = currentAddress.ID
			createdAt = currentAddress.CreatedAt
		} else {
			addressId = uuid.NewString()
			createdAt = time.Now().UTC()
		}

		newAddress := addressDao.CreateAddressDataParams{
			ID:           &addressId,
			City:         apiResult.City,
			Neighborhood: apiResult.Neighborhood(),
			State:        apiResult.State(),
			Longitude:    apiResult.Longitude,
			Latitude:     apiResult.Latitude,
			CreatedAt:    &createdAt,
			ProfileID:    profileId,
		}

		return saveAddress(ctx, newAddress)
	}
}
