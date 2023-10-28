package models

import "time"

type AttachmentSchema struct {
	ID          string    `json:"id"`
	ResourceURL string    `json:"resourceUrl"`
	CreatedAt   time.Time `json:"createdAt"`
}
