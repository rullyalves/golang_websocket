package models

import (
	"time"
)

type ChatView struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	MatchID        string    `json:"matchId"`
	ParticipantIDs []string  `json:"participantIds"`
}
