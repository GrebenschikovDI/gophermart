package usecase

import (
	"context"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
	"github.com/pkg/errors"
)

type BalanceUseCase struct {
	balanceRepo repository.BalanceRepository
}

func NewBalanceUseCase(balanceRepo repository.BalanceRepository) *BalanceUseCase {
	return &BalanceUseCase{
		balanceRepo: balanceRepo,
	}
}

func (u *BalanceUseCase) Get(ctx context.Context, userID int) (*entity.Balance, error) {
	balance, err := u.balanceRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting balance of uesr: %d", userID)
	}
	return balance, nil
}
func (u *BalanceUseCase) Add(ctx context.Context, userID int, amount float64) error {
	err := u.balanceRepo.Add(ctx, userID, amount)
	if err != nil {
		return errors.Wrapf(err, "can't add balance to user: %d", userID)
	}
	return nil
}

func (u *BalanceUseCase) Withdraw(ctx context.Context, userID int, withdraw float64) error {
	tx, err := u.balanceRepo.BeginTransaction()
	if err != nil {
		return errors.Wrapf(err, "wihtdraw transaction error, user: %d", userID)
	}
	defer tx.Rollback()
	canWithdraw, err := u.balanceRepo.Check(ctx, tx, userID, withdraw)
	if err != nil {
		return errors.Wrapf(err, "wihtdraw check error, user: %d", userID)
	}
	if !canWithdraw {
		return errors.Wrapf(gophermart.ErrLowBalance, "low balance, user: %d", userID)
	}
	err = u.balanceRepo.Withdraw(ctx, tx, userID, withdraw)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return errors.Wrapf(err, "wihtdraw commit error, user: %d", userID)
	}
	return nil
}
