package models

import "time"

type MatchView struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}
