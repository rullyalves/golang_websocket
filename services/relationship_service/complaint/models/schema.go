package models

import (
	profile "go_ws/services/profile_service/profile/models"
	category "go_ws/services/relationship_service/complaint_category/models"
	"time"
)

type ComplaintSchema struct {
	ID          string                            `json:"id"`
	Description string                            `json:"description"`
	CreatedAt   time.Time                         `json:"createdAt"`
	Category    *category.ComplaintCategorySchema `json:"category"`
	Sender      *profile.ProfileSchema            `json:"sender"`
	Receiver    *profile.ProfileSchema            `json:"receiver"`
}
