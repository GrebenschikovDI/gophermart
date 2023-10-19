package repository

import (
	"context"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
)

type BalanceRepository interface {
	Create(ctx context.Context, balance *entity.Balance) error
	GetByID(ctx context.Context, id int) (*entity.Balance, error)
	GetByUserID(ctx context.Context, userID int) ([]*entity.Balance, error)
	Update(ctx context.Context, balance *entity.Balance) error
}
