package input

type CreateDislikeParamsInput struct {
	SenderID   *string `json:"senderId,omitempty"`
	ReceiverID string  `json:"receiverId"`
}
