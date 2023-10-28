package models

import "time"

type ConfigurationSchema struct {
	ID                    string    `json:"id"`
	MinimumAppBuildNumber int       `json:"minimumAppBuildNumber"`
	CreatedAt             time.Time `json:"createdAt"`
}
