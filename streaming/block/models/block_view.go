package models

import "time"

type BlockView struct {
	ID         string    `json:"id" bson:"id"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	MatchID    string    `json:"matchId" bson:"matchId"`
	ReceiverID string    `json:"receiverId" bson:"receiverId"`
}
