package models

import (
	profile "go_ws/services/profile_service/profile/models"
	complaint "go_ws/services/relationship_service/complaint/models"
	"time"
)

type BlockSchema struct {
	ID        *string                    `json:"id,omitempty"`
	CreatedAt time.Time                  `json:"createdAt"`
	Complaint *complaint.ComplaintSchema `json:"complaint,omitempty"`
	Sender    *profile.ProfileSchema     `json:"sender"`
	Receiver  *profile.ProfileSchema     `json:"receiver"`
}
