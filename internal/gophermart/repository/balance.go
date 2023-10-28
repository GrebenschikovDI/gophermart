package repository

import (
	"context"
	"database/sql"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
)

type BalanceRepository interface {
	Create(ctx context.Context, balance *entity.Balance) error
	GetByID(ctx context.Context, userID int) (*entity.Balance, error)
	Add(ctx context.Context, userID int, amount float64) error
	BeginTransaction() (*sql.Tx, error)
	Check(ctx context.Context, tx *sql.Tx, userID int, withdrawal float64) (bool, error)
	Withdraw(ctx context.Context, tx *sql.Tx, userID int, withdraw float64) error
}
