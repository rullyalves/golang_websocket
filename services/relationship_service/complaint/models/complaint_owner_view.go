package models

type ComplaintOwnerView struct {
	OwnerID   string        `json:"owner_id"`
	Complaint ComplaintView `json:"data"`
}
