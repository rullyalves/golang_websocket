package usecases

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	imageDao "go_ws/services/profile_service/image/dao/neo4j"
	optionDao "go_ws/services/profile_service/options/dao/neo4j"
	profileDao "go_ws/services/profile_service/profile/dao/neo4j"
	"go_ws/shared/collections"
	shared "go_ws/shared/models"
	"go_ws/shared/queue"
	streamingModels "go_ws/streaming/shared/models"
	"time"
)

type UpdateProfileParams struct {
	Name        *string   `json:"name" validate:"omitempty,min=2"`
	Age         *int      `json:"age" validate:"omitempty,min=18"`
	Height      *float32  `json:"height" validate:"omitempty,min=1"`
	Description *string   `json:"description" validate:"omitempty,min=5"`
	OptionIds   *[]string `json:"optionIds" validate:"omitempty,dive,uuid4"`
	ImageUrls   *[]string `json:"imageUrls" validate:"omitempty,dive,url"`
}

type UpdateProfileData func(context.Context, string, UpdateProfileParams) error

func NewUpdateProfile(neo4jDriver *neo4j.DriverWithContext, sqsClient *sqs.Client) UpdateProfileData {
	saveProfile := profileDao.Save(neo4jDriver)
	saveImages := imageDao.LinkWithProfile(neo4jDriver)
	saveOptions := optionDao.LinkWithProfile(neo4jDriver)

	sendToQueue := queue.SendMessageToQueue(sqsClient)
	updatePreferences := UpdateProfile(saveProfile, sendToQueue, saveImages, saveOptions)
	return updatePreferences
}

func UpdateProfile(saveProfile profileDao.SaveProfile, sendToQueue queue.SendToQueue, linkImages imageDao.LinkImageWithProfile, linkOptions optionDao.LinkOptionsWithProfile) UpdateProfileData {
	return func(ctx context.Context, profileId string, params UpdateProfileParams) error {

		newProfile := profileDao.CreateProfileDataParams{
			ID:          profileId,
			Name:        params.Name,
			Age:         params.Age,
			Height:      params.Height,
			Description: params.Description,
		}

		err := saveProfile(ctx, newProfile)

		if err != nil {
			return err
		}

		optionIds := params.OptionIds

		if optionIds != nil {
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

		//TODO: enviar evento para todo mundo que deu like / tem match com esse perfil
		profileEvent := streamingModels.CreateEventPayload[any](
			//TODO: pegar dados do dislike
			nil,
			"",
			time.Now().UTC(),
			shared.MessageTypeUpdate,
		)

		return sendToQueue(ctx, queue.ProfileStreamQueueUrl, []queue.MessageInput{profileEvent})
	}
}
