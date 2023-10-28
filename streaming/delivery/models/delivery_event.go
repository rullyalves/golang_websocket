package models

import (
	shared "go_ws/shared/models"
	"time"
)

type DeliveryEvent struct {
	ID        string             `json:"id" json:"id"`
	Data      DeliveryView       `json:"data" json:"data"`
	UserID    string             `json:"userId" bson:"userId"`
	Type      shared.MessageType `json:"type" bson:"type"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
