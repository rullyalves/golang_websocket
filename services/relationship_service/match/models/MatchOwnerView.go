package models

type MatchOwnerView struct {
	OwnerID string    `json:"owner_id"`
	Match   MatchView `json:"data"`
}
