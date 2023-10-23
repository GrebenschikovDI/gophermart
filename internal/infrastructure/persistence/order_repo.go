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
	_, err := r.db.ExecContext(ctx, "INSERT INTO orders (id, user_id, status) values ($1, $2, $3)", order.ID,
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

func (r orderRepo) Update(ctx context.Context, id int, status string) (*entity.Order, error) {
	_, err := r.db.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r orderRepo) GetByUserID(ctx context.Context, userID int) ([]*entity.Order, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, status, uploaded_at FROM orders WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*entity.Order, 0)
	for rows.Next() {
		order := &entity.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
