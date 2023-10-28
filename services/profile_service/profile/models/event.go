package models

import shared "go_ws/shared/models"

type ProfileSubscriptionEvent struct {
	ID          string             `json:"id"`
	MessageType shared.MessageType `json:"messageType"`
	Data        *ProfileSchema     `json:"data"`
}
