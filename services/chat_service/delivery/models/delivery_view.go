package models

import "time"

// TODO: mudar forma como isso é salvo / consultado no banco
type DeliveryStatus string

const (
	DeliveryStatusSent      DeliveryStatus = "sent"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusRead      DeliveryStatus = "read"
	DeliveryStatusPlayed    DeliveryStatus = "played"
)

type DeliveryView struct {
	ID        string         `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	Status    DeliveryStatus `json:"status"`
	TargetID  string         `json:"targetId"`
	MessageID string         `json:"messageId"`
}
