package models

import "time"

type LocationSchema struct {
	ID           string    `json:"id"`
	City         string    `json:"city"`
	Neighborhood string    `json:"neighborhood"`
	State        string    `json:"state"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	CreatedAt    time.Time `json:"createdAt"`
}
