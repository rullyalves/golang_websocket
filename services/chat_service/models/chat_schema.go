package models

import (
	profile "go_ws/services/profile_service/profile/models"
	"time"
)

type ChatSchema struct {
	ID           string                   `json:"id"`
	CreatedAt    time.Time                `json:"createdAt"`
	Participants []*profile.ProfileSchema `json:"participants"`
	Messages     []MessageSchema          `json:"messages"`
}
