package models

import (
	profile "go_ws/services/profile_service/profile/models"
	"time"
)

type DislikeSchema struct {
	ID        string                 `json:"id"`
	CreatedAt time.Time              `json:"createdAt"`
	Sender    *profile.ProfileSchema `json:"sender"`
	Receiver  *profile.ProfileSchema `json:"receiver"`
}
