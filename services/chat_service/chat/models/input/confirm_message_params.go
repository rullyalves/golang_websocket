package input

import "time"

type ConfirmMessageReadParams struct {
	date      time.Time `json:"readAt"`
	profileId string    `json:"profileId"`
	messageId string    `json:"messageId"`
}
