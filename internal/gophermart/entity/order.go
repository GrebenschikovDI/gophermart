package entity

import "time"

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}
