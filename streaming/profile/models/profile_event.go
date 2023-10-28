package models

import (
	"go_ws/shared/models"
	"time"
)

type ProfileEvent struct {
	ID        string             `json:"id" bson:"id"`
	Data      ProfileView        `json:"data" bson:"data"`
	UserId    string             `json:"userId" bson:"userId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	Type      models.MessageType `json:"type" bson:"type"`
}
