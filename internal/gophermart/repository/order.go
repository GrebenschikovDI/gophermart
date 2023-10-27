package repository

import (
	"context"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id string) (*entity.Order, error)
	Update(ctx context.Context, id, status string, accrual *float64) (*entity.Order, error)
	GetByUserID(ctx context.Context, userID int) ([]*entity.Order, error)
	GetToSend(ctx context.Context, offset, limit int) ([]*entity.Order, error)
}
