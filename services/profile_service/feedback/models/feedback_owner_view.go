package models

type FeedbackOwnerView struct {
	OwnerID  string       `json:"ownerId"`
	Feedback FeedbackView `json:"data"`
}
