package input

type CreateFeedbackParamsInput struct {
	Description *string `json:"description,omitempty"`
	ProfileID   *string `json:"profileId,omitempty"`
}
