package models

import (
	sharedModels "go_ws/services/chat_service/models"
	"time"
)

type MessageView struct {
	ID        string                 `json:"id" bson:"id"`
	Text      *string                `json:"text" bson:"text"`
	MediaType sharedModels.MediaType `json:"mediaType" bson:"mediaType"`
	CreatedAt time.Time              `json:"createdAt" bson:"createdAt"`
	ChatID    string                 `json:"chatId" bson:"chatId"`
	SenderID  string                 `json:"senderId" bson:"senderId"`
	ParentID  *string                `json:"parentId" bson:"parent"`
}
