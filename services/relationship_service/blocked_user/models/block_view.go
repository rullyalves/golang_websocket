package models

import "time"

type BlockView struct {
	ID        string    `json:"id"`
	BlockedID string    `json:"blockedId"`
	CreatedAt time.Time `json:"createdAt"`
}
