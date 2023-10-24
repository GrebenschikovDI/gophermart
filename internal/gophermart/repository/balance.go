package repository

import (
	"context"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
)

type BalanceRepository interface {
	Create(ctx context.Context, balance *entity.Balance) error
	GetByID(ctx context.Context, userID int) (*entity.Balance, error)
	Add(ctx context.Context, userID int, amount float64) error
	Withdraw(ctx context.Context, userID int, withdraw float64) error
	CheckWithdrawal(ctx context.Context, userID int, withdrawal float64) (bool, error)
}
