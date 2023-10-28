package models

import "time"

type PreferencesView struct {
	ID        string    `json:"id"`
	MinAge    int       `json:"min_age"`
	MaxAge    int       `json:"max_age"`
	Distance  int       `json:"distance"`
	CreatedAt time.Time `json:"createdAt"`
}
