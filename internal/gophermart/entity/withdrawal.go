package entity

import "time"

type Withdrawal struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	OrderID     int       `json:"order_id"`
	Amount      float64   `json:"amount"`
	Total       float64   `json:"total"`
	ProcessedAt time.Time `json:"processed_at"`
}
