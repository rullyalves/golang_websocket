package models

type ProfileOwnerView struct {
	OwnerID string      `json:"owner_id"`
	Profile ProfileView `json:"data"`
}
