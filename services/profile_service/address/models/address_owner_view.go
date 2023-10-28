package models

type AddressOwnerView struct {
	OwnerID string      `json:"ownerId"`
	Address AddressView `json:"data"`
}
