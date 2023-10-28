package models

import (
	"go_ws/shared/models"
	"time"
)

type BlockEvent struct {
	ID        string             `json:"id" bson:"id"`
	Data      BlockView          `json:"data" bson:"data"`
	UserID    string             `json:"userId" bson:"userId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	Type      models.MessageType `json:"type" bson:"type"`
}
