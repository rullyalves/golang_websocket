package models

import "time"

type Settings struct {
	ID                             string    `json:"id"`
	IsVisible                      bool      `json:"isVisible"`
	AllowReceiveMatchNotifications bool      `json:"allowReceiveMatchNotifications"`
	AllowReceiveLikeNotifications  bool      `json:"allowReceiveLikeNotifications"`
	CreatedAt                      time.Time `json:"createdAt"`
}
