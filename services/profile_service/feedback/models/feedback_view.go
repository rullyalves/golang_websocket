package models

import "time"

type FeedbackView struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UserID      string    `json:"user_id"`
}
