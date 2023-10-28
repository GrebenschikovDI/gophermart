package entity

import "time"

type Withdrawal struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	OrderID     string    `json:"order_id"`
	Amount      float64   `json:"amount"`
	ProcessedAt time.Time `json:"processed_at"`
}
