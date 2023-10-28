package models

import "time"

type ImageView struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}
