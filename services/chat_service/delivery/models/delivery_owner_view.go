package models

type DeliveryOwnerView struct {
	OwnerID  string       `json:"ownerId"`
	Delivery DeliveryView `json:"data"`
}
