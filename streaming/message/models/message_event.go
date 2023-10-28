package models

import (
	"go_ws/shared/models"
	"time"
)

type MessageEvent struct {
	ID        string             `json:"id" bson:"id"`
	Data      MessageView        `json:"data" bson:"data"`
	UserId    string             `json:"userId" bson:"userId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	Type      models.MessageType `json:"type" bson:"type"`
}
