package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	imageDao "go_ws/services/profile_service/image/dao/neo4j"
	optionDao "go_ws/services/profile_service/options/dao/neo4j"
	preferencesDao "go_ws/services/profile_service/preferences/dao/neo4j"
	profileDao "go_ws/services/profile_service/profile/dao/neo4j"
	"go_ws/shared/collections"
	"time"
)

type CreateProfileParams struct {
	Name        *string    `json:"name" validate:"omitempty,min=2"`
	Age         *int       `json:"age" validate:"omitempty,min=18"`
	Height      *float32   `json:"height" validate:"omitempty,min=1"`
	Description *string    `json:"description" validate:"omitempty,min=5"`
	CreatedAt   *time.Time `json:"createdAt" validate:"omitempty,datetime"`
	OptionIds   *[]string  `json:"optionIds" validate:"omitempty,dive,uuid4"`
	ImageUrls   *[]string  `json:"imageUrls" validate:"omitempty,dive,url"`
}

type CreateProfileData func(context.Context, CreateProfileParams) error

func NewCreateProfile(neo4jDriver *neo4j.DriverWithContext) CreateProfileData {
	saveProfile := profileDao.Save(neo4jDriver)
	savePreferences := preferencesDao.Save(neo4jDriver)
	saveImages := imageDao.LinkWithProfile(neo4jDriver)
	saveOptions := optionDao.LinkWithProfile(neo4jDriver)

	updatePreferences := CreateProfile(saveProfile, savePreferences, saveImages, saveOptions)
	return updatePreferences
}

func CreateProfile(saveProfile profileDao.SaveProfile, savePreferences preferencesDao.SavePreferences, linkImages imageDao.LinkImageWithProfile, linkOptions optionDao.LinkOptionsWithProfile) CreateProfileData {
	return func(ctx context.Context, params CreateProfileParams) error {

		//TODO: criar user

		profileId := uuid.NewString()
		newProfile := profileDao.CreateProfileDataParams{
			ID:          profileId,
			Name:        params.Name,
			Age:         params.Age,
			Height:      params.Height,
			Description: params.Description,
			CreatedAt:   params.CreatedAt,
		}

		err := saveProfile(ctx, newProfile)

		if err != nil {
			return err
		}

		minAge := params.Age
		maxAge := *params.Age + 6
		minDistance := 10 * 1000
		createdAt := time.Now().UTC()

		if *minAge >= 22 {
			*minAge = (*minAge) - 4
		}

		err = savePreferences(ctx, profileId, preferencesDao.CreatePreferencesDataParams{
			ID:        uuid.NewString(),
			Distance:  &minDistance,
			MinAge:    minAge,
			MaxAge:    &maxAge,
			CreatedAt: &createdAt,
		})

		if err != nil {
			return err
		}

		optionIds := params.OptionIds

		if optionIds != nil && len(*optionIds) > 0 {
			err := linkOptions(ctx, profileId, *optionIds)
			if err != nil {
				return err
			}
		}

		imageUrls := params.ImageUrls

		if imageUrls != nil && len(*imageUrls) > 0 {
			images := collections.Map(*imageUrls, func(url string) imageDao.CreateImageDataParams {
				return imageDao.CreateImageDataParams{
					ID:        uuid.NewString(),
					Url:       url,
					CreatedAt: time.Now(),
				}
			})
			err := linkImages(ctx, profileId, images)
			if err != nil {
				return err
			}
		}

		return err
	}
}
