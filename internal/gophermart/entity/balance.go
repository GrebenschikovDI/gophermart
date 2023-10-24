package entity

import "time"

type Balance struct {
	UserID      int       `json:"user_id"`
	Amount      float64   `json:"amount"`
	Withdrawn   float64   `json:"withdrawn"`
	ProcessedAt time.Time `json:"processed_at"`
}
