package usecase

import (
	"context"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
	"github.com/pkg/errors"
)

type WithdrawalUseCase struct {
	withdrawalRepository repository.WithdrawalRepository
}

func NewWithdrawalUseCase(withdrawalRepository repository.WithdrawalRepository) *WithdrawalUseCase {
	return &WithdrawalUseCase{
		withdrawalRepository: withdrawalRepository,
	}
}

func (w *WithdrawalUseCase) Add(ctx context.Context, withdrawal *entity.Withdrawal) error {
	err := w.withdrawalRepository.Create(ctx, withdrawal)
	if err != nil {
		return errors.Wrapf(err, "creating withdrawal err, userID: %d", withdrawal.UserID)
	}
	return nil
}

func (w *WithdrawalUseCase) GetList(ctx context.Context, userID int) ([]*entity.Withdrawal, error) {
	withdrawals, err := w.withdrawalRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "getting withdrawal list err, userID: %d", userID)
	}
	return withdrawals, nil
}
