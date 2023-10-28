package models

import (
	profile "go_ws/services/profile_service/profile/models"
	"time"
)

type DeliverySchema struct {
	ID         string                 `json:"id"`
	CreatedAt  time.Time              `json:"createdAt"`
	ReadAt     *time.Time             `json:"readAt,omitempty"`
	PlayedAt   *time.Time             `json:"playedAt,omitempty"`
	DeliveryAt time.Time              `json:"deliveryAt"`
	Receiver   *profile.ProfileSchema `json:"receiver"`
	//Message    *message.MessageSchema `json:"message"`
}
