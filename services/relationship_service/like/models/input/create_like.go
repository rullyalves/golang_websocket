package input

type CreateLikeParamsInput struct {
	SenderID   *string `json:"senderId,omitempty"`
	ReceiverID string  `json:"receiverId"`
}
