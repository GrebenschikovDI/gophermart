package entity

import "time"

type Balance struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	OrderID     int       `json:"order_id"`
	Amount      float64   `json:"amount"`
	ProcessedAt time.Time `json:"processed_at"`
}
