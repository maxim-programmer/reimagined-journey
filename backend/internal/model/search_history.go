package model

import "time"

type SearchHistory struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
}