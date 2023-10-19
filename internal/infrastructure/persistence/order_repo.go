package persistence

import (
	"context"
	"database/sql"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/entity"
	"github.com/GrebenschikovDI/gophermart.git/internal/gophermart/repository"
)

type orderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) repository.OrderRepository {
	return &orderRepo{
		db: db,
	}
}

func (r orderRepo) Create(ctx context.Context, order *entity.Order) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO orders (user_id, status) values ($1, $2)",
		order.UserID, order.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r orderRepo) GetByID(ctx context.Context, id int) (*entity.Order, error) {
	row := r.db.QueryRowContext(ctx, "SELECT  id, user_id, status, uploaded_at FROM orders WHERE id = $1", id)
	order := &entity.Order{}
	err := row.Scan(&order.ID, &order.UserID, &order.Status, &order.UploadedAt)
	if err != nil {
		return nil, err
	}
	return order, nil
}
