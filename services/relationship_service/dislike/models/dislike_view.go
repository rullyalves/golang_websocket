package models

import "time"

type DislikeView struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
