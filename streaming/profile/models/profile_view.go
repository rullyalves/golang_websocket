package models

import "time"

type ProfileView struct {
	ID          string    `json:"id" bson:"id"`
	Name        string    `json:"name" bson:"name"`
	Age         int       `json:"age" bson:"age"`
	Height      *int      `json:"height" bson:"height"`
	Description string    `json:"description" bson:"description"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}
