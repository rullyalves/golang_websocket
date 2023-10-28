package models

import (
	"go_ws/services/chat_service/delivery/models"
	"time"
)

type DeliveryView struct {
	ID        string                `json:"id" bson:"id"`
	CreatedAt time.Time             `json:"createdAt" bson:"createdAt"`
	Status    models.DeliveryStatus `json:"status" bson:"status"`
	TargetID  string                `json:"targetId" bson:"targetId"`
	MessageId string                `json:"messageId" bson:"messageId"`
}
