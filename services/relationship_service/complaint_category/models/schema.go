package models

import "time"

type ComplaintCategorySchema struct {
	ID        string                     `json:"id"`
	Name      string                     `json:"name"`
	CreatedAt time.Time                  `json:"createdAt"`
	Children  []*ComplaintCategorySchema `json:"children"`
}
