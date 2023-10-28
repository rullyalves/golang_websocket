package models

type ChatOwnerView struct {
	OwnerID    string         `json:"ownerId"`
	Attachment AttachmentView `json:"data"`
}
