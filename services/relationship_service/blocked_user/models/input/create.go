package input

type CreateBlockParamsInput struct {
	SenderID   *string `json:"senderId,omitempty"`
	ReceiverID string  `json:"receiverId"`
}
