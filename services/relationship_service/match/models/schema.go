package models

import "time"

type MatchSchema struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	IsActive  bool      `json:"isActive"`
	//Chat      *ChatSchema `json:"chat"`
}
