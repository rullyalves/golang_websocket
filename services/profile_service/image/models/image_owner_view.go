package models

type ImageOwnerView struct {
	OwnerID string    `json:"ownerId"`
	Image   ImageView `json:"data"`
}
