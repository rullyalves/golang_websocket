package models

import "time"

type ChatView struct {
	ID           string    `json:"id" bson:"id"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	MatchID      string    `json:"matchId" bson:"matchId"`
	Participants []string  `json:"participants" bson:"participants"`
}
