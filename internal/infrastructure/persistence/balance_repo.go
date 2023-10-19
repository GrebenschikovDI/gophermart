package persistence

import (
	"context"
	"database/sql"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
)

type balanceRepo struct {
	db *sql.DB
}

func NewBalanceRepo(db *sql.DB) repository.BalanceRepository {
	return &balanceRepo{
		db: db,
	}
}

func (r balanceRepo) Create(ctx context.Context, balance *entity.Balance) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO balance (user_id, order_id, amount) VALUES ($1, $2, $3)",
		balance.UserID, balance.OrderID, balance.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (r balanceRepo) GetByID(ctx context.Context, id int) (*entity.Balance, error) {
	row := r.db.QueryRowContext(ctx, "SELECT  id, user_id, order_id, amount, processed_at FROM balance WHERE id = $1", id)
	balance := &entity.Balance{}
	err := row.Scan(&balance.ID, &balance.UserID, &balance.OrderID, &balance.Amount, &balance.ProcessedAt)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (r balanceRepo) GetByUserID(ctx context.Context, userID int) ([]*entity.Balance, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, order_id, amount, processed_at FROM balance WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []*entity.Balance
	for rows.Next() {
		balance := &entity.Balance{}
		err := rows.Scan(&balance.ID, &balance.UserID, &balance.OrderID, &balance.Amount, &balance.ProcessedAt)
		if err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}
	return balances, err
}

func (r balanceRepo) Update(ctx context.Context, balance *entity.Balance) error {
	//TODO implement me
	panic("implement me")
}
