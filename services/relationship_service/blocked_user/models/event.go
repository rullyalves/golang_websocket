package models

import shared "go_ws/shared/models"

type BlockSubscriptionEvent struct {
	ID          string             `json:"id"`
	MessageType shared.MessageType `json:"messageType"`
	Data        *BlockSchema       `json:"data"`
}
