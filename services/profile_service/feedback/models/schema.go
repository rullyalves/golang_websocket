package models

import "time"

type FeedbackSchema struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	//Profile     *ProfileSchema `json:"profile"`
}
