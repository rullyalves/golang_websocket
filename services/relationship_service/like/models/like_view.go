package models

import "time"

type LikeView struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	SenderID   string    `json:"senderId"`
	ReceiverID string    `json:"receiverId"`
}
