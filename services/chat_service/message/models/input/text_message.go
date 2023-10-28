package input

import "time"

type CreateMessageParamsInput struct {
	ID              string    `json:"id" validate:"required,uuid4"`
	CreatedAt       time.Time `json:"createdAt" validate:"required"`
	SenderID        string    `json:"senderId" validate:"required,uuid4"`
	ChatID          string    `json:"chatId" validate:"required,uuid4"`
	ParentMessageID *string   `json:"parentMessageId,omitempty" validate:"omitempty,required,uuid4"`
}

type CreateTextMessageParamsInput struct {
	Text    string                    `json:"text" validate:"required"`
	Message *CreateMessageParamsInput `json:"message" validate:"required"`
}
