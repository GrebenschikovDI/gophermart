package repository

import (
	"context"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id int) (*entity.Order, error)
}
