package models

import (
	profile "go_ws/services/profile_service/profile/models"
	"time"
)

type MessageSchema struct {
	ID        string                 `json:"id"`
	CreatedAt time.Time              `json:"createdAt"`
	SentAt    time.Time              `json:"sentAt"`
	Text      *string                `json:"text,omitempty"`
	MediaType MediaType              `json:"mediaType"`
	Sender    *profile.ProfileSchema `json:"sender"`
	Chat      *ChatSchema            `json:"chat"`
	Parent    *MessageSchema         `json:"parent,omitempty"`
}
