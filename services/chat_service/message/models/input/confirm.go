package input

import "time"

type ConfirmMessageDeliveryParamsInput struct {
	DeliveryAt time.Time `json:"deliveryAt"`
	ProfileID  string    `json:"profileId"`
	MessageID  string    `json:"messageId"`
}

type ConfirmMessagePlayedParamsInput struct {
	PlayedAt  time.Time `json:"playedAt"`
	ProfileID string    `json:"profileId"`
	MessageID string    `json:"messageId"`
}

type ConfirmMessageReadParamsInput struct {
	ReadAt    time.Time `json:"readAt"`
	ProfileID string    `json:"profileId"`
	MessageID string    `json:"messageId"`
}
