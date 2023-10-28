package models

import (
	shared "go_ws/shared/models"
	"time"
)

type ChatEvent struct {
	Data      ChatView           `json:"data"`
	ID        string             `json:"id"`
	UserID    string             `json:"userId"`
	CreatedAt time.Time          `json:"createdAt"`
	Type      shared.MessageType `json:"type"`
}
