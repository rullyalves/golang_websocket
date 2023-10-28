package models

import "time"

type WebRtcRoom struct {
	RoomID     string                       `json:"roomId" bson:"roomId"`
	ChatID     string                       `json:"chatId" bson:"chatId"`
	Offer      *RtcSessionDescription       `json:"offer" bson:"offer"`
	Answer     []RtcSessionDescription      `json:"answer" bson:"answer"`
	Candidates map[string][]RtcIceCandidate `json:"candidates" bson:"candidates"`
	CreatedAt  time.Time                    `json:"createdAt" bson:"createdAt"`
}
