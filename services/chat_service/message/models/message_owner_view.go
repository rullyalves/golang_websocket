package models

import "go_ws/streaming/message/models"

type MessageOwnerView struct {
	OwnerId string             `json:"ownerId"`
	Message models.MessageView `json:"data"`
}
