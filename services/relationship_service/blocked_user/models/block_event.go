package models

import (
	shared "go_ws/shared/models"
	"time"
)

type BlockEvent struct {
	Data      BlockView          `json:"data"`
	ID        string             `json:"id"`
	UserId    string             `json:"userId"`
	CreatedAt time.Time          `json:"createdAt"`
	Type      shared.MessageType `json:"type"`
}
