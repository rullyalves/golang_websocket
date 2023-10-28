package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	filterDao "go_ws/services/profile_service/filter/dao/neo4j"
	options "go_ws/services/profile_service/options/models"
	preferencesDao "go_ws/services/profile_service/preferences/dao/neo4j"
	"go_ws/shared/collections"
	"time"
)

type MultipleOption struct {
	Type       options.OptionType `json:"type"`
	IsRequired bool               `json:"isRequired"`
	Data       []string           `json:"data"`
}

type UpdatePreferencesData struct {
	MinAge   *int              `json:"min_age" validate:"omitempty,min=18"`
	MaxAge   *int              `json:"max_age" validate:"omitempty,min=18"`
	Distance *int              `json:"distance" validate:"omitempty,min=1000"`
	Filters  *[]MultipleOption `json:"filters"`
}

type UpdateProfilePreferences func(ctx context.Context, profileId string, params UpdatePreferencesData) error

func NewUpdatePreferences(neo4jDriver *neo4j.DriverWithContext) UpdateProfilePreferences {
	savePreferences := preferencesDao.Save(neo4jDriver)
	updateFilters := filterDao.SaveAll(neo4jDriver)
	updatePreferences := UpdatePreferences(savePreferences, updateFilters)
	return updatePreferences
}

func UpdatePreferences(savePreferences preferencesDao.SavePreferences, updateFilters filterDao.SaveAllFilters) UpdateProfilePreferences {
	return func(ctx context.Context, preferencesId string, params UpdatePreferencesData) error {

		newPreferences := preferencesDao.CreatePreferencesDataParams{
			ID:       preferencesId,
			MinAge:   params.MinAge,
			MaxAge:   params.MaxAge,
			Distance: params.Distance,
		}

		err := savePreferences(ctx, preferencesId, newPreferences)

		if err != nil {
			return err
		}

		filters := params.Filters

		if filters == nil || len(*filters) == 0 {
			return nil
		}

		newOptions := collections.Map(*params.Filters, func(value MultipleOption) filterDao.CreateFilterDataParams {
			return filterDao.CreateFilterDataParams{
				ID:         uuid.NewString(),
				OptionIds:  value.Data,
				IsRequired: value.IsRequired,
				Type:       value.Type,
				CreatedAt:  time.Now().UTC(),
			}
		})

		return updateFilters(ctx, preferencesId, newOptions)
	}
}
