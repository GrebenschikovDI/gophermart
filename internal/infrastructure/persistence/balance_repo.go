package persistence

import (
	"context"
	"database/sql"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
	"time"
)

type balanceRepo struct {
	db *sql.DB
}

func NewBalanceRepo(db *sql.DB) repository.BalanceRepository {
	return &balanceRepo{
		db: db,
	}
}

func (b *balanceRepo) Create(ctx context.Context, balance *entity.Balance) error {
	_, err := b.db.ExecContext(ctx, "INSERT INTO balance (user_id,  amount) VALUES ($1, $2)",
		balance.UserID, balance.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (b *balanceRepo) GetByID(ctx context.Context, userID int) (*entity.Balance, error) {
	row := b.db.QueryRowContext(
		ctx,
		"SELECT user_id, amount, withdrawn, processed_at FROM balance WHERE user_id = $1",
		userID)
	balance := &entity.Balance{}
	err := row.Scan(&balance.UserID, &balance.Amount, &balance.Withdrawn, &balance.ProcessedAt)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (b *balanceRepo) Add(ctx context.Context, userID int, amount float64) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	t := time.Now().Format(time.RFC3339)

	query := `
		INSERT INTO balance (user_id, amount, processed_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE 
		SET amount = balance.amount + excluded.amount, processed_at = excluded.processed_at
	`
	_, err = tx.ExecContext(ctx, query, userID, amount, t)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *balanceRepo) BeginTransaction() (*sql.Tx, error) {
	return b.db.Begin()
}

func (b *balanceRepo) Check(ctx context.Context, tx *sql.Tx, userID int, withdrawal float64) (bool, error) {
	var currentAmount float64
	query := "SELECT amount FROM balance WHERE user_id = $1 FOR UPDATE"
	err := tx.QueryRowContext(ctx, query, userID).Scan(&currentAmount)
	if err != nil {
		return false, err
	}
	if withdrawal > currentAmount {
		return false, nil
	}
	return true, nil
}

func (b *balanceRepo) Withdraw(ctx context.Context, tx *sql.Tx, userID int, withdraw float64) error {
	t := time.Now().Format(time.RFC3339)

	query := `
		INSERT INTO balance (user_id, withdrawn, processed_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE 
		SET amount = balance.amount - excluded.withdrawn, withdrawn = balance.withdrawn + excluded.withdrawn, processed_at = excluded.processed_at
	`
	_, err := tx.ExecContext(ctx, query, userID, withdraw, t)
	if err != nil {
		return err
	}
	return nil
}
