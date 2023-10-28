package models

import (
	"go_ws/services/profile_service/options/models"
	"time"
)

type FilterView struct {
	ID            string            `json:"id"`
	OptionIds     []string          `json:"optionIds"`
	IsRequired    bool              `json:"isRequired"`
	Type          models.OptionType `json:"type"`
	PreferencesID string            `json:"preferencesId"`
	CreatedAt     time.Time         `json:"createdAt"`
}
