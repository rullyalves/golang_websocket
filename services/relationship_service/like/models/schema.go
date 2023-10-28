package models

import "time"

type LikeSchema struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	//Sender    *ProfileSchema `json:"sender"`
	//Receiver  *ProfileSchema `json:"receiveR"`
}
