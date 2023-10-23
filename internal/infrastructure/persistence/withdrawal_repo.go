package persistence

import (
	"context"
	"database/sql"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
)

type withdrawalRepo struct {
	db *sql.DB
}

func NewWithdrawalRepo(db *sql.DB) repository.WithdrawalRepository {
	return &withdrawalRepo{
		db: db,
	}
}

func (r withdrawalRepo) Create(ctx context.Context, withdrawal *entity.Withdrawal) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO withdrawals (user_id, order_id, amount, total) VALUES ($1, $2, $3)",
		withdrawal.UserID, withdrawal.OrderID, withdrawal.Amount, withdrawal.Total)
	if err != nil {
		return err
	}
	return nil
}

func (r withdrawalRepo) GetByID(ctx context.Context, id int) (*entity.Withdrawal, error) {
	row := r.db.QueryRowContext(ctx, "SELECT  id, user_id, order_id, amount, total, processed_at FROM withdrawals WHERE id = $1", id)
	withdrawal := &entity.Withdrawal{}
	err := row.Scan(&withdrawal.ID, &withdrawal.UserID, &withdrawal.OrderID, &withdrawal.Amount, &withdrawal.ProcessedAt)
	if err != nil {
		return nil, err
	}
	return withdrawal, nil
}

func (r withdrawalRepo) GetByUserID(ctx context.Context, userID int) ([]*entity.Withdrawal, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, order_id, amount, total, processed_at FROM withdrawals WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []*entity.Withdrawal
	for rows.Next() {
		withdrawal := &entity.Withdrawal{}
		err := rows.Scan(&withdrawal.ID, &withdrawal.UserID, &withdrawal.OrderID, &withdrawal.Amount, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
