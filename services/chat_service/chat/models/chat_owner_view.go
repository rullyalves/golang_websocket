package models

type ChatOwnerView struct {
	OwnerID string   `json:"ownerId"`
	Chat    ChatView `json:"data"`
}
