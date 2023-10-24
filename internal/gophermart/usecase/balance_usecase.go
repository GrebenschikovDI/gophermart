package usecase

import (
	"context"
	"errors"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
)

var ErrLowBalance = errors.New("balance is low")

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
		return nil, err
	}
	return balance, nil
}
func (u *BalanceUseCase) Add(ctx context.Context, userID int, amount float64) error {
	err := u.balanceRepo.Add(ctx, userID, amount)
	if err != nil {
		return err
	}
	return nil
}
func (u *BalanceUseCase) Withdraw(ctx context.Context, userID int, withdraw float64) error {
	isValid, err := u.balanceRepo.CheckWithdrawal(ctx, userID, withdraw)
	if err != nil {
		return err
	}
	if isValid {
		err := u.balanceRepo.Withdraw(ctx, userID, withdraw)
		if err != nil {
			return err
		}
	} else {
		return ErrLowBalance
	}
	return nil
}
