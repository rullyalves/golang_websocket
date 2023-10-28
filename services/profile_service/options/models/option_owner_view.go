package models

type OptionOwnerView struct {
	OwnerId string     `json:"ownerId"`
	Data    OptionView `json:"data"`
}
