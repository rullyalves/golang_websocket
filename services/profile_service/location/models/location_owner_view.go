package models

type LocationOwnerView struct {
	OwnerID  string       `json:"ownerId"`
	Location LocationView `json:"data"`
}
