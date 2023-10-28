package models

import (
	shared "go_ws/shared/models"
	"time"
)

type ChatEvent struct {
	ID        string             `json:"id" bson:"id"`
	Data      ChatView           `json:"data" bson:"data"`
	UserID    string             `json:"userId" bson:"userId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	Type      shared.MessageType `json:"type" bson:"type"`
}
