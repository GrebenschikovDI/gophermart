package repository

import (
	"context"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
)

type WithdrawalRepository interface {
	Create(ctx context.Context, withdrawal *entity.Withdrawal) error
	GetByID(ctx context.Context, id int) (*entity.Withdrawal, error)
	GetByUserID(ctx context.Context, userID int) ([]*entity.Withdrawal, error)
}
