package models

import (
	"time"
)

type DislikeView struct {
	ID        string    `json:"id" bson:"id"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	SenderID  string    `json:"senderId" bson:"senderId"`
}
